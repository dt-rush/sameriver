package sameriver

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	r := NewRateLimiter(10 * time.Millisecond)
	x := 0
	n := 128
	for i := 0; i < n; i++ {
		r.Do(func() {
			x += 1
		})
		time.Sleep(5 * time.Millisecond)
	}
	if x > n/2 {
		t.Fatal("did not rate limit")
	}
}

func TestRateLimiterReset(t *testing.T) {
	r := NewRateLimiter(10 * time.Millisecond)
	x := 0
	r.Do(func() {
		x = 1
	})
	r.Reset()
	r.Do(func() {
		x = 2
	})
	if x != 2 {
		t.Fatal("did not reset rate limiter")
	}
}
