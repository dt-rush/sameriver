// Array-Based Read-Write Queueing Lock
// implementation of a circular array-based queueing lock which is
// safe from overflow at the expense of relatively-busy checking on a single
// atomic value (dequeueCount) in the case of overflow
package engine

import (
	"sync/atomic"
	"time"
)

// This lock is used for entire-component locks, and hence the queue size is
// tuned for that purpose - we assume that when a whole component is locked, for
// example during physics / spatial hash computation / draw / collision, we
// may accumulate about 256 lock attempts on that component for various entities
// during the time it is locked, as an upper bound.
const COMPONENT_ACCESS_ABQL_QUEUE_SZ = 256

// active modification lock sleeps for 5 us each check
// this means that when a whole-component lock unlocks, each of the
// goroutines trying to lock on this component in lockEntityComponent will
// dequeue as readers around every 5 microseconds, so we have a minimum time of
// (256 * 5 us) which it will take for the goroutines to all get their component,
// or 1.28 milliseconds (leaves plenty of time in the update loop for them to
// do their action and release the component before the next time physics /
// spatial hash / draw / collision will need to lock the entity component
// by locking this lock as a writer
const COMPONENT_ACCESS_LOCK_SLEEP = 5 * time.Microsecond

type ComponentAccessLock struct {
	arr          []int
	ticket       uint32
	dequeueCount uint32
	nReaders     uint32
}

func NewComponentAccessLock() *ComponentAccessLock {
	abql := ComponentAccessLock{
		arr:          make([]int, queue_sz),
		ticket:       0,
		dequeueCount: 0,
		nReaders:     0}
	abql.arr[0] = 1
	return &abql
}

func (l *ComponentAccessLock) RLock() {
	ticket := atomic.AddUint32(&l.ticket, 1) - 1
	for ticket-l.dequeueCount >= uint32(COMPONENT_ACCESS_LOCK_QUEUE_SZ) {
		time.Sleep(COMPONENT_ACCESS_LOCK_SLEEP_NS)
	}
	for l.arr[ticket%uint32(COMPONENT_ACCESS_LOCK_QUEUE_SZ)] != 1 {
		time.Sleep(COMPONENT_ACCESS_LOCK_SLEEP_NS)
	}
	// increment nReaders
	atomic.AddUint32(&l.nReaders, 1)
	// move the queue forward after incrementing nReaders, so that
	// either another call to RLock() can get the queue head or else
	// a call to Lock() can get the queue head but wait for nReaders = 0
	l.arr[int(ticket)%COMPONENT_ACCESS_LOCK_QUEUE_SZ] = 0
	l.arr[int(ticket+1)%COMPONENT_ACCESS_LOCK_QUEUE_SZ] = 1
	atomic.AddUint32(&l.dequeueCount, 1)
	return
}

func (l *ComponentAccessLock) Lock() uint32 {
	ticket := atomic.AddUint32(&l.ticket, 1) - 1
	for ticket-l.dequeueCount >= uint32(COMPONENT_ACCESS_LOCK_QUEUE_SZ) {
		time.Sleep(COMPONENT_ACCESS_LOCK_SLEEP_NS)
	}
	for l.arr[ticket%uint32(COMPONENT_ACCESS_LOCK_QUEUE_SZ)] != 1 {
		time.Sleep(COMPONENT_ACCESS_LOCK_SLEEP_NS)
	}
	// wait here if our turn in the queue came because the prior lock was
	// an RLock releasing itself in hopes of triggering another RLock() instance
	for l.nReaders > 0 {
		time.Sleep(COMPONENT_ACCESS_LOCK_SLEEP_NS)
	}
	return ticket
}

func (l *ComponentAccessLock) RUnlock() {
	// decrement nReaders
	atomic.AddUint32(&l.nReaders, ^uint32(0))
}

func (l *ComponentAccessLock) Unlock(ticket uint32) {
	l.arr[int(ticket)%COMPONENT_ACCESS_LOCK_QUEUE_SZ] = 0
	l.arr[int(ticket+1)%COMPONENT_ACCESS_LOCK_QUEUE_SZ] = 1
	atomic.AddUint32(&l.dequeueCount, 1)
}
