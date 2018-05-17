package engine

type SignalSendRateLimiter struct {
	mutex sync.Mutex
	guard sync.Once
	out   chan (bool)
	delay time.Duration
}

func (r *SignalRateLimiter) Do(signal bool) {
	r.mutex.Lock()
	r.guard.Do(func() {
		r.out <- signal
		go func() {
			time.Sleep(r.delay)
			r.mutex.Lock()
			r.guard = sync.Once{}
			r.mutex.Unlock()
		}()
	})
	r.mutex.Unlock()
}
