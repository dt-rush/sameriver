package sameriver

import (
	"time"

	"go.uber.org/atomic"
)

// Like the above, but it can be reset while sleeping
type RateLimiter struct {
	limited atomic.Uint32
	delay   time.Duration
}

func NewRateLimiter(delay time.Duration) *RateLimiter {
	return &RateLimiter{delay: delay}
}

func (r *RateLimiter) Do(f func()) {
	if r.limited.CompareAndSwap(0, 1) {
		f()
		go func() {
			time.Sleep(r.delay)
			r.limited.CompareAndSwap(1, 0)
		}()
	}
}

func (r *RateLimiter) Reset() {
	r.limited.CompareAndSwap(1, 0)
}

func (r *RateLimiter) Limited() bool {
	return r.limited.Load() == 1
}
