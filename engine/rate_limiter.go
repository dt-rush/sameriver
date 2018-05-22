package engine

import (
	"sync"
	"sync/atomic"
	"time"
)

// Basically a once with a delay after which the once becomes usable again
type RateLimiter struct {
	once  sync.Once
	mutex sync.RWMutex
	delay time.Duration
}

func (r *RateLimiter) Do(f func()) {
	r.mutex.RLock()
	r.once.Do(f)
	r.mutex.RUnlock()

	time.Sleep(r.delay)

	r.mutex.Lock()
	r.once = sync.Once{}
	r.mutex.Unlock()
}

// Like the above, but it can be reset while sleeping
type ResettableRateLimiter struct {
	once  sync.Once
	mutex sync.RWMutex
	delay time.Duration
	// used so the automatic reset can check, after sleeping, if
	// another goroutine also had called Reset() while it slept. If so,
	// do not reset as we would if nothing happened during sleep.
	resetCounter uint32
}

func (r *ResettableRateLimiter) Do(f func()) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	r.once.Do(func() {
		f()
		go func() {
			resetCounterBeforeSleep := atomic.LoadUint32(&r.resetCounter)
			time.Sleep(r.delay)
			if atomic.LoadUint32(&r.resetCounter) == resetCounterBeforeSleep {
				r.Reset()
			}
		}()
	})
}

func (r *ResettableRateLimiter) Reset() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	atomic.AddUint32(&r.resetCounter, 1)
	r.once = sync.Once{}
}
