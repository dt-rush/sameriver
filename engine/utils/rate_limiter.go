package utils

import (
	"go.uber.org/atomic"
	"sync"
	"time"
)

// Basically a once with a delay after which the once becomes usable again
type RateLimiter struct {
	once  sync.Once
	mutex sync.RWMutex
	delay time.Duration
}

func NewRateLimiter(delay time.Duration) *RateLimiter {
	return &RateLimiter{delay: delay}
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
	resetCounter atomic.Uint32
}

func NewResettableRateLimiter(delay time.Duration) *ResettableRateLimiter {
	return &ResettableRateLimiter{delay: delay}
}

func (r *ResettableRateLimiter) Do(f func()) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	r.once.Do(func() {
		f()
		go func() {
			resetCounterBeforeSleep := r.resetCounter.Load()
			time.Sleep(r.delay)
			if r.resetCounter.Load() == resetCounterBeforeSleep {
				r.Reset()
			}
		}()
	})
}

func (r *ResettableRateLimiter) Reset() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.resetCounter.Inc()
	r.once = sync.Once{}
}
