package ratelimiter

import (
	"sync"
	"time"
)

type ResetRateLimiter struct {
	maxRequests int
	timeWindow  time.Duration
	timestamps  []time.Time
	mu          sync.Mutex
}

func NewRateLimiterReset(maxRequests int, timeWindow time.Duration) *ResetRateLimiter {
	return &ResetRateLimiter{
		maxRequests: maxRequests,
		timeWindow:  timeWindow,
		timestamps:  make([]time.Time, 0),
		mu:          sync.Mutex{},
	}
}

func (rl *ResetRateLimiter) getResetTime() time.Time {
	now := time.Now()
	// Calculate how much time has passed in the current window
	elapsed := now.UnixNano() % int64(rl.timeWindow)
	// Subtracting the elapsed time from now will give the start time of the current window
	startOfCurrentWindow := now.Add(-time.Duration(elapsed))
	return startOfCurrentWindow.Add(rl.timeWindow)
}

func (rl *ResetRateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Reset the timestamps if the first entry is older than the last reset time
	if len(rl.timestamps) > 0 && rl.timestamps[0].Before(rl.getResetTime().Add(-rl.timeWindow)) {
		rl.timestamps = []time.Time{}
	}

	if len(rl.timestamps) < rl.maxRequests {
		rl.timestamps = append(rl.timestamps, time.Now())
		return true
	}
	return false
}
