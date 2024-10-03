package clients

import (
	"sync"
	"time"

	"studentgit.kata.academy/quant/torque/pkg/ratelimiter"

	"studentgit.kata.academy/quant/torque/internal/adapters/clients/types/okx"
)

const (
	maxRequests60 = 60
	maxRequests10 = 10
	timeInterval2 = time.Second * 2
)

// RateLimiterMap is a concurrent-safe map for rate limiters.
type RateLimiterMap struct {
	mu          sync.RWMutex
	internalMap map[okx.Endpoint]*ratelimiter.RateLimiter
}

// NewRateLimiterMap creates a new concurrent-safe rate limiter map.
func NewRateLimiterMap() *RateLimiterMap {
	return &RateLimiterMap{
		//nolint:exhaustive // We don't need to rate limit other endpoints
		internalMap: map[okx.Endpoint]*ratelimiter.RateLimiter{
			okx.Order:       ratelimiter.NewRateLimiter(maxRequests60, timeInterval2),
			okx.CancelOrder: ratelimiter.NewRateLimiter(maxRequests60, timeInterval2),
			okx.AmendOrder:  ratelimiter.NewRateLimiter(maxRequests60, timeInterval2),
			okx.Balances:    ratelimiter.NewRateLimiter(maxRequests10, timeInterval2),
		},
		mu: sync.RWMutex{},
	}
}

// Get returns the rate limiter for the given endpoint and ensures that the access is concurrent-safe.
func (m *RateLimiterMap) Get(endpoint okx.Endpoint) *ratelimiter.RateLimiter {
	m.mu.RLock()
	defer m.mu.RUnlock()
	rateLimiter := m.internalMap[endpoint]
	return rateLimiter
}

// Set sets the rate limiter for the given endpoint and ensures that the access is concurrent-safe.
func (m *RateLimiterMap) Set(endpoint okx.Endpoint, limiter *ratelimiter.RateLimiter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.internalMap[endpoint] = limiter
}

// Delete removes the rate limiter for the given endpoint and ensures that the access is concurrent-safe.
func (m *RateLimiterMap) Delete(endpoint okx.Endpoint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.internalMap, endpoint)
}
