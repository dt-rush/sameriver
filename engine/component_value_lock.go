// Array-Based Read-Write Queueing Lock
// implementation of a circular array-based queueing lock which is
// safe from overflow at the expense of relatively-busy checking on a single
// atomic value (dequeueCount) in the case of overflow
package engine

import (
	"sync/atomic"
	"time"
)

// This lock is used by AtomicEntityModify and SafeGet(), hence the constants
// are tuned for those purposes

// at any given time, we can expect that there may be around upper bound 4
// simultaneous reads / writes wanting a component value
// TODO: make all such variables as these tunable in a single location
// so that different games with different requirements can set up what makes
// sense for them
const COMPONENT_VALUE_LOCK_QUEUE_SZ = 4

// active modification lock sleeps for 10 us each check
const COMPONENT_VALUE_LOCK_SLEEP = 10 * time.Microsecond

type ComponentValueLock struct {
	arr          []int
	ticket       uint32
	dequeueCount uint32
	nReaders     uint32
}

func NewComponentValueLock() *ComponentValueLock {
	abql := ComponentValueLock{
		arr:          make([]int, queue_sz),
		ticket:       0,
		dequeueCount: 0,
		nReaders:     0}
	abql.arr[0] = 1
	return &abql
}

func (l *ComponentValueLock) RLock() {
	ticket := atomic.AddUint32(&l.ticket, 1) - 1
	for ticket-l.dequeueCount >= uint32(COMPONENT_VALUE_LOCK_QUEUE_SZ) {
		time.Sleep(COMPONENT_VALUE_LOCK_SLEEP)
	}
	for l.arr[ticket%uint32(COMPONENT_VALUE_LOCK_QUEUE_SZ)] != 1 {
		time.Sleep(COMPONENT_VALUE_LOCK_SLEEP)
	}
	// increment nReaders
	atomic.AddUint32(&l.nReaders, 1)
	// move the queue forward after incrementing nReaders, so that
	// either another call to RLock() can get the queue head or else
	// a call to Lock() can get the queue head but wait for nReaders = 0
	l.arr[int(ticket)%COMPONENT_VALUE_LOCK_QUEUE_SZ] = 0
	l.arr[int(ticket+1)%COMPONENT_VALUE_LOCK_QUEUE_SZ] = 1
	atomic.AddUint32(&l.dequeueCount, 1)
	return
}

func (l *ComponentValueLock) Lock() uint32 {
	ticket := atomic.AddUint32(&l.ticket, 1) - 1
	for ticket-l.dequeueCount >= uint32(COMPONENT_VALUE_LOCK_QUEUE_SZ) {
		time.Sleep(COMPONENT_VALUE_LOCK_SLEEP)
	}
	for l.arr[ticket%uint32(COMPONENT_VALUE_LOCK_QUEUE_SZ)] != 1 {
		time.Sleep(COMPONENT_VALUE_LOCK_SLEEP)
	}
	// wait here if our turn in the queue came because the prior lock was
	// an RLock releasing itself in hopes of triggering another RLock() instance
	for l.nReaders > 0 {
		time.Sleep(COMPONENT_VALUE_LOCK_SLEEP)
	}
	return ticket
}

func (l *ComponentValueLock) RUnlock() {
	// decrement nReaders
	atomic.AddUint32(&l.nReaders, ^uint32(0))
}

func (l *ComponentValueLock) Unlock(ticket uint32) {
	l.arr[int(ticket)%COMPONENT_VALUE_LOCK_QUEUE_SZ] = 0
	l.arr[int(ticket+1)%COMPONENT_VALUE_LOCK_QUEUE_SZ] = 1
	atomic.AddUint32(&l.dequeueCount, 1)
}
