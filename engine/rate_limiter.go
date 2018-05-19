package engine

import (
	"sync"
	"time"
)

type RateLimiter struct {
	once  sync.Once
	mutex sync.RWMutex
	delay time.Duration
}

func (r *RateLimiter) Invoke(f func()) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	f()
	time.Sleep(r.delay)
	r.once = sync.Once{}
}
