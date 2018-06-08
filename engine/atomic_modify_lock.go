// Array-Based Read-Write Queueing Lock
// implementation of a circular array-based queueing lock which is
// safe from overflow at the expense of relatively-busy checking on a single
// atomic value (dequeueCount) in the case of overflow
package engine

import (
	"sync/atomic"
	"time"
)

// This lock is used by AtomicEntityModify, hence the constants are
// tuned for those purposes

// active modification lock holds an array queue of 4 (can overflow, but really
// 4 simultaneous modifiers is a good upper bound, considering that at even in
// some crazy situation, an entity, two other entities, and an entity class could
// desire to modify the entity, or an entity, another entity, and two world logic
// goroutines, etc.)
const ATOMIC_MODIFY_LOCK_QUEUE_SZ = 4

// active modification lock sleeps for 10 us each check
const ATOMIC_MODIFY_LOCK_SLEEP = 10 * time.Microsecond

type AtomicModifyLock struct {
	arr          []int
	ticket       uint32
	dequeueCount uint32
	nReaders     uint32
}

func NewAtomicModifyLock() *AtomicModifyLock {
	abql := AtomicModifyLock{
		arr:          make([]int, queue_sz),
		ticket:       0,
		dequeueCount: 0,
		nReaders:     0}
	abql.arr[0] = 1
	return &abql
}

func (l *AtomicModifyLock) RLock() {
	ticket := atomic.AddUint32(&l.ticket, 1) - 1
	for ticket-l.dequeueCount >= uint32(ATOMIC_MODIFY_LOCK_QUEUE_SZ) {
		time.Sleep(ATOMIC_MODIFY_LOCK_SLEEP)
	}
	for l.arr[ticket%uint32(ATOMIC_MODIFY_LOCK_QUEUE_SZ)] != 1 {
		time.Sleep(ATOMIC_MODIFY_LOCK_SLEEP)
	}
	// increment nReaders
	atomic.AddUint32(&l.nReaders, 1)
	// move the queue forward after incrementing nReaders, so that
	// either another call to RLock() can get the queue head or else
	// a call to Lock() can get the queue head but wait for nReaders = 0
	l.arr[int(ticket)%ATOMIC_MODIFY_LOCK_QUEUE_SZ] = 0
	l.arr[int(ticket+1)%ATOMIC_MODIFY_LOCK_QUEUE_SZ] = 1
	atomic.AddUint32(&l.dequeueCount, 1)
	return
}

func (l *AtomicModifyLock) Lock() uint32 {
	ticket := atomic.AddUint32(&l.ticket, 1) - 1
	for ticket-l.dequeueCount >= uint32(ATOMIC_MODIFY_LOCK_QUEUE_SZ) {
		time.Sleep(ATOMIC_MODIFY_LOCK_SLEEP)
	}
	for l.arr[ticket%uint32(ATOMIC_MODIFY_LOCK_QUEUE_SZ)] != 1 {
		time.Sleep(ATOMIC_MODIFY_LOCK_SLEEP)
	}
	// wait here if our turn in the queue came because the prior lock was
	// an RLock releasing itself in hopes of triggering another RLock() instance
	for l.nReaders > 0 {
		time.Sleep(ATOMIC_MODIFY_LOCK_SLEEP)
	}
	return ticket
}

func (l *AtomicModifyLock) RUnlock() {
	// decrement nReaders
	atomic.AddUint32(&l.nReaders, ^uint32(0))
}

func (l *AtomicModifyLock) Unlock(ticket uint32) {
	l.arr[int(ticket)%ATOMIC_MODIFY_LOCK_QUEUE_SZ] = 0
	l.arr[int(ticket+1)%ATOMIC_MODIFY_LOCK_QUEUE_SZ] = 1
	atomic.AddUint32(&l.dequeueCount, 1)
}
