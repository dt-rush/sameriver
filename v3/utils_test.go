package sameriver

import (
	"github.com/dt-rush/sameriver/v3/utils"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	r := utils.NewRateLimiter(50 * time.Millisecond)
	x := 0
	for i := 0; i < 16; i++ {
		r.Do(func() {
			x += 1
		})
	}
	if x == 16 {
		t.Fatal("did not rate limit")
	}
}
