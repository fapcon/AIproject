package cache

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	parent *Cache[K, V]
	mu     sync.RWMutex
	data   map[K]V
	key    func(V) K
}

func NewCache[K comparable, V any](key func(V) K) *Cache[K, V] {
	return &Cache[K, V]{
		parent: nil,
		mu:     sync.RWMutex{},
		data:   make(map[K]V, 0),
		key:    key,
	}
}

func (c *Cache[K, V]) Reset(values ...V) {
	if c.parent != nil {
		panic("Reset is not supported in transaction")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[K]V, len(values))

	for _, value := range values {
		key := c.key(value)
		c.data[key] = value
	}
}

func (c *Cache[K, V]) Add(values ...V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, value := range values {
		key := c.key(value)
		c.data[key] = value
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	if ok || c.parent == nil {
		return value, ok
	}

	return c.parent.Get(key)
}

func (c *Cache[K, V]) GetAll() []V {
	var res []V
	c.iterate(func(value V) bool {
		res = append(res, value)
		return true
	})
	return res
}

func (c *Cache[K, V]) Filter(condition func(V) bool) []V {
	var res []V
	c.iterate(func(value V) bool {
		if condition(value) {
			res = append(res, value)
		}
		return true
	})
	return res
}

func (c *Cache[K, V]) FindFirst(condition func(V) bool) (V, bool) {
	var res V
	var found bool
	c.iterate(func(value V) bool {
		if condition(value) {
			res = value
			found = true
			return false // stop iteration
		}
		return true
	})
	return res, found
}

func (c *Cache[K, V]) iterate(process func(V) bool) {
	count := 0
	last := 0
	c1 := c
	for c1 != nil {
		c1.mu.RLock()
		defer c1.mu.RUnlock()
		last = len(c1.data)
		count += last
		c1 = c1.parent
	}

	keys := make(map[K]struct{}, count-last)

	c2 := c
	for c2.parent != nil {
		for key, value := range c2.data {
			if _, ok := keys[key]; !ok {
				if !process(value) {
					return
				}
				keys[key] = struct{}{}
			}
		}
		c2 = c2.parent
	}

	// Following code is outside the loop to optimise memory (and time).
	// The last cache doesn't have parent, so we don't need to update keys map.
	for key, value := range c2.data {
		if _, ok := keys[key]; !ok {
			if !process(value) {
				return
			}
		}
	}
}
