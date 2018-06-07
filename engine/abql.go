package engine

import (
	"sync/atomic"
	"time"
)

// implementation of a circular array-based queueing lock which is
// safe from overflow at the expense of relatively-busy checking on a single
// atomic value (dequeueCount) in the case of overflow
type ABQL struct {
	arr          [ABQL_QUEUE_SZ]int
	sleep        time.Duration
	ticket       uint32
	dequeueCount uint32
}

func NewABQL(sleep time.Duration) *ABQL {
	abql := ABQL{sleep: sleep}
	abql.arr[0] = 1
	return &abql
}

func (l *ABQL) Lock() int {
	ticket := atomic.AddUint32(&l.ticket, 1) - 1
	for ticket-atomic.LoadUint32(&l.dequeueCount) >= uint32(ABQL_QUEUE_SZ) {
		time.Sleep(l.sleep)
	}
	for l.arr[ticket%uint32(ABQL_QUEUE_SZ)] != 1 {
		time.Sleep(l.sleep)
	}
	return int(ticket)
}

func (l *ABQL) Unlock(ticket int) {
	l.arr[ticket%ABQL_QUEUE_SZ] = 0
	l.arr[(ticket+1)%ABQL_QUEUE_SZ] = 1
	atomic.AddUint32(&l.dequeueCount, 1)
}
