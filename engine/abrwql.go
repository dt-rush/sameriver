package rwpql

import (
	"go.uber.org/atomic"
	"time"
)

// Array-Based Read-Write Queueing Lock
type ABRWQL struct {
	arr      []int
	queue_sz int

	sleep_ns     time.Duration
	ticket       *atomic.Uint32
	dequeueCount *atomic.Uint32
	nReaders     *atomic.Uint32
}

func NewABRWQL(queue_sz int, sleep_ns int) *ABRWQL {
	abql := ABRWQL{
		arr:          make([]int, queue_sz),
		queue_sz:     queue_sz,
		sleep_ns:     time.Duration(sleep_ns),
		ticket:       atomic.NewUint32(0),
		dequeueCount: atomic.NewUint32(0),
		nReaders:     atomic.NewUint32(0)}
	abql.arr[0] = 1
	return &abql
}

func (l *ABRWQL) RLock() {
	ticket := l.ticket.Inc() - 1
	for ticket-l.dequeueCount.Load() >= uint32(l.queue_sz) {
		time.Sleep(l.sleep_ns)
	}
	for l.arr[ticket%uint32(l.queue_sz)] != 1 {
		time.Sleep(l.sleep_ns)
	}
	l.nReaders.Inc()
	// move the queue forward after incrementing nReaders, so that
	// either another call to RLock() can get the queue head or else
	// a call to Lock() can get the queue head but wait for nReaders = 0
	l.arr[int(ticket)%l.queue_sz] = 0
	l.arr[int(ticket+1)%l.queue_sz] = 1
	l.dequeueCount.Inc()
	return
}

func (l *ABRWQL) Lock() uint32 {
	ticket := l.ticket.Inc() - 1
	for ticket-l.dequeueCount.Load() >= uint32(l.queue_sz) {
		time.Sleep(l.sleep_ns)
	}
	for l.arr[ticket%uint32(l.queue_sz)] != 1 {
		time.Sleep(l.sleep_ns)
	}
	// wait here if our turn in the queue came because the prior lock was
	// an RLock releasing itself in hopes of triggering another RLock() instance
	for l.nReaders.Load() > 0 {
		time.Sleep(l.sleep_ns)
	}
	return ticket
}

func (l *ABRWQL) RUnlock() {
	readersRemaining := l.nReaders.Dec()
}

func (l *ABRWQL) Unlock(ticket uint32) {
	l.arr[int(ticket)%l.queue_sz] = 0
	l.arr[int(ticket+1)%l.queue_sz] = 1
	l.dequeueCount.Inc()
}
