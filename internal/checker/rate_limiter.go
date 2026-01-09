package checker

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu          sync.Mutex
	lastRequest time.Time
	minInterval time.Duration
}

func NewRateLimiter(reqPerSecond int) *RateLimiter {
	return &RateLimiter{
		minInterval: time.Second / time.Duration(reqPerSecond),
	}
}

func (rl *RateLimiter) Wait() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	elapsed := time.Since(rl.lastRequest)
	if elapsed < rl.minInterval {
		time.Sleep(rl.minInterval - elapsed)
	}
	rl.lastRequest = time.Now()
}
