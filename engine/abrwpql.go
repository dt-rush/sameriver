package rwpql

import (
	"go.uber.org/atomic"
	"time"
)

// Array-Based Read-Write-Priority Queueing Lock
type ABRWPQL struct {
	arr      []int
	queue_sz int

	sleep_ns        time.Duration
	lockState       *atomic.Uint32
	ticket          *atomic.Uint32
	dequeueCount    *atomic.Uint32
	nReaders        *atomic.Uint32
	priorityWaiting *atomic.Uint32
}

func NewABRWPQL(queue_sz int, sleep_ns int) *ABRWPQL {
	abql := ABRWPQL{
		arr:             make([]int, queue_sz),
		queue_sz:        queue_sz,
		sleep_ns:        time.Duration(sleep_ns),
		lockState:       atomic.NewUint32(OPEN),
		ticket:          atomic.NewUint32(0),
		dequeueCount:    atomic.NewUint32(0),
		nReaders:        atomic.NewUint32(0),
		priorityWaiting: atomic.NewUint32(0)}
	abql.arr[0] = 1
	return &abql
}

func (l *ABRWPQL) yieldToPriorityLock() {
	l.lockState.Store(PRIORITY_LOCK_RESERVED)
	for l.lockState.Load() != OPEN {
		time.Sleep(l.sleep_ns)
	}
}

func (l *ABRWPQL) RLock() {
	ticket := l.ticket.Inc() - 1
	for ticket-l.dequeueCount.Load() >= uint32(l.queue_sz) {
		time.Sleep(l.sleep_ns)
	}
	for l.arr[ticket%uint32(l.queue_sz)] != 1 {
		time.Sleep(l.sleep_ns)
	}
	for {
		for !(l.lockState.Load() == RLOCKED || l.lockState.CAS(OPEN, RLOCKED)) {
			time.Sleep(l.sleep_ns)
		}
		if l.priorityWaiting.Load() > 0 {
			l.yieldToPriorityLock()
		} else {
			l.nReaders.Inc()
			// move the queue forward after incrementing nReaders, so that
			// either another call to RLock() can get the queue head or else
			// a call to Lock() can get the queue head but wait for nReaders = 0
			l.arr[int(ticket)%l.queue_sz] = 0
			l.arr[int(ticket+1)%l.queue_sz] = 1
			l.dequeueCount.Inc()
			return
		}
	}
}

func (l *ABRWPQL) Lock() uint32 {
	ticket := l.ticket.Inc() - 1
	for ticket-l.dequeueCount.Load() >= uint32(l.queue_sz) {
		time.Sleep(l.sleep_ns)
	}
	for l.arr[ticket%uint32(l.queue_sz)] != 1 {
		time.Sleep(l.sleep_ns)
	}
	for l.nReaders.Load() > 0 {
		time.Sleep(l.sleep_ns)
	}
	for {
		for !l.lockState.CAS(OPEN, LOCKED) {
			time.Sleep(l.sleep_ns)
		}
		if l.priorityWaiting.Load() > 0 {
			l.yieldToPriorityLock()
		} else {
			return ticket
		}
	}
}

func (l *ABRWPQL) PLock() {
	l.priorityWaiting.Inc()
	for !(l.lockState.CAS(OPEN, PRIORITY_LOCKED) ||
		l.lockState.CAS(PRIORITY_LOCK_RESERVED, PRIORITY_LOCKED)) {
		time.Sleep(l.sleep_ns)
	}
	l.priorityWaiting.Dec()
}

func (l *ABRWPQL) RUnlock() {
	readersRemaining := l.nReaders.Dec()
	if readersRemaining == 0 {
		l.lockState.Store(OPEN)
	}
}

func (l *ABRWPQL) Unlock(ticket uint32) {
	l.arr[int(ticket)%l.queue_sz] = 0
	l.arr[int(ticket+1)%l.queue_sz] = 1
	l.lockState.Store(OPEN)
	l.dequeueCount.Inc()
}

func (l *ABRWPQL) PUnlock() {
	l.lockState.Store(OPEN)
}
