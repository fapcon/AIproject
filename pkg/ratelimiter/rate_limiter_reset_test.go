package ratelimiter_test

import (
	"testing"
	"time"

	"studentgit.kata.academy/quant/torque/pkg/ratelimiter"
)

func TestRateLimiter_Allow(t *testing.T) {
	maxRequests := 5
	timeWindow := 10 * time.Second
	limiter := ratelimiter.NewRateLimiterReset(maxRequests, timeWindow)

	// Test: Should allow 5 requests in the 10-second window
	for i := 0; i < maxRequests; i++ {
		if !limiter.Allow() {
			t.Fatalf("expected Allow() to return true for request %d", i+1)
		}
	}

	// Test: Shouldn't allow a 6th request in the same 10-second window
	if limiter.Allow() {
		t.Fatal("expected Allow() to return false after maxRequests in the same time window")
	}

	// Test: After waiting for the 10-second interval, it should allow new requests
	time.Sleep(timeWindow)
	if !limiter.Allow() {
		t.Fatal("expected Allow() to return true after waiting for the time window to pass")
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	maxRequests := 2
	timeWindow := 2 * time.Second
	limiter := ratelimiter.NewRateLimiterReset(maxRequests, timeWindow)

	// Test: Should allow 2 requests in the 2-second window
	for i := 0; i < maxRequests; i++ {
		if !limiter.Allow() {
			t.Fatalf("expected Allow() to return true for request %d", i+1)
		}
	}

	// Wait until the start of the next time window to make the next request.
	// This ensures that we are not dependent on the precise timing
	// of the test execution relative to the rate limiter's internal clock.
	time.Sleep(timeWindow)

	// Test: After the time window has fully elapsed, the next request should be allowed
	if !limiter.Allow() {
		t.Fatal("expected Allow() to return true after the time window has fully reset")
	}

	// Immediately try another request, which should be allowed because we're in a new window.
	if !limiter.Allow() {
		t.Fatal("expected Allow() to return true for the second request in a new time window")
	}

	// Attempt a third request in the same window, which should be rejected.
	if limiter.Allow() {
		t.Fatal("expected Allow() to return false for the third request in the same time window")
	}
}

func TestRateLimiter_ResetOnExactInterval(t *testing.T) {
	maxRequests := 2
	timeWindow := 5 * time.Second
	limiter := ratelimiter.NewRateLimiterReset(maxRequests, timeWindow)

	// Wait for the next exact 5-second interval (like 00, 05, 10, etc.)
	targetTime := time.Now().Truncate(timeWindow).Add(timeWindow)
	time.Sleep(time.Until(targetTime))

	// Now, make 2 requests
	for i := 0; i < maxRequests; i++ {
		if !limiter.Allow() {
			t.Fatalf("expected Allow() to return true for request %d", i+1)
		}
	}

	// Wait for the next exact 5-second interval
	targetTime = time.Now().Truncate(timeWindow).Add(timeWindow)
	time.Sleep(time.Until(targetTime))

	// It should reset and allow new requests again
	if !limiter.Allow() {
		t.Fatal("expected Allow() to return true after waiting for the exact interval to reset")
	}
}
