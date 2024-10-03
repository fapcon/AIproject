package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	maxRequests int
	timeWindow  time.Duration
	timestamps  []time.Time
	mu          sync.Mutex
}

func NewRateLimiter(maxRequests int, timeWindow time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		timeWindow:  timeWindow,
		timestamps:  make([]time.Time, 0),
		mu:          sync.Mutex{},
	}
}

func (rl *RateLimiter) slideWindow() {
	cutoffTime := time.Now().Add(-rl.timeWindow)
	idx := 0
	for idx < len(rl.timestamps) && rl.timestamps[idx].Before(cutoffTime) {
		idx++
	}
	rl.timestamps = rl.timestamps[idx:]
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.slideWindow()
	if len(rl.timestamps) < rl.maxRequests {
		rl.timestamps = append(rl.timestamps, time.Now())
		return true
	}
	return false
}
