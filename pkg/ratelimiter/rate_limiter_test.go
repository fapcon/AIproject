package ratelimiter_test

import (
	"testing"
	"time"

	ratelimiter "studentgit.kata.academy/quant/torque/pkg/ratelimiter"
)

func TestRateLimiter(t *testing.T) {
	limiter := ratelimiter.NewRateLimiter(5, 1*time.Second)

	// Should allow the next 5 requests
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Fatalf("expected request %d to be allowed, but it was not", i+1)
		}
	}

	// The 6th request should not be allowed
	if limiter.Allow() {
		t.Fatal("expected request 6 to be denied, but it was allowed")
	}

	// After 1 second, another request should be allowed
	time.Sleep(1 * time.Second)
	if !limiter.Allow() {
		t.Fatal("expected request after 1 second to be allowed, but it was not")
	}
}
