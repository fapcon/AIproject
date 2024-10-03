package cache

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Deletable[K comparable, V any] struct {
	parent  *Deletable[K, V]
	mu      sync.RWMutex
	data    map[K]V
	removed map[K]struct{}
	key     func(V) K
}

func NewDeletable[K comparable, V any](key func(V) K) *Deletable[K, V] {
	return &Deletable[K, V]{
		parent:  nil,
		mu:      sync.RWMutex{},
		data:    make(map[K]V, 0),
		removed: nil,
		key:     key,
	}
}

func (c *Deletable[K, V]) Reset(values ...V) {
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

func (c *Deletable[K, V]) Add(values ...V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, value := range values {
		key := c.key(value)
		c.data[key] = value
		delete(c.removed, key)
	}
}

func (c *Deletable[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	if ok || c.parent == nil {
		return value, ok
	}

	if _, removed := c.removed[key]; removed {
		return value, false
	}

	return c.parent.Get(key)
}

func (c *Deletable[K, V]) GetAll() []V {
	var res []V
	c.iterate(func(value V) bool {
		res = append(res, value)
		return true
	})
	return res
}

func (c *Deletable[K, V]) Filter(condition func(V) bool) []V {
	var res []V
	c.iterate(func(value V) bool {
		if condition(value) {
			res = append(res, value)
		}
		return true
	})
	return res
}

func (c *Deletable[K, V]) FindFirst(condition func(V) bool) (V, bool) {
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

func (c *Deletable[K, V]) Delete(keys ...K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.parent != nil {
		for _, key := range keys {
			c.removed[key] = struct{}{}
			delete(c.data, key)
		}
		return
	}

	for _, key := range keys {
		delete(c.data, key)
	}
}

func (c *Deletable[K, V]) iterate(process func(V) bool) { //nolint:gocognit // don't know how to simplify
	countRemoved := 0
	countData := 0
	lastData := 0 // there is no lastRemoved because it would be always zero (because rootCache.removed == nil)
	c1 := c
	for c1 != nil {
		c1.mu.RLock()
		defer c1.mu.RUnlock()
		countRemoved += len(c1.removed)
		lastData = len(c1.data)
		countData += lastData
		c1 = c1.parent
	}

	removed := make(map[K]struct{}, countRemoved)
	keys := make(map[K]struct{}, countData-lastData)

	c2 := c
	for c2.parent != nil {
		for key, value := range c2.data {
			if _, ok := removed[key]; ok {
				continue
			}
			if _, ok := keys[key]; ok {
				continue
			}
			if !process(value) {
				return
			}
			keys[key] = struct{}{}
		}
		for key := range c2.removed {
			removed[key] = struct{}{}
		}
		c2 = c2.parent
	}

	// Following code is outside the loop to optimise memory (and time).
	// The last cache doesn't have parent, so we don't need to update keys and removed maps.
	for key, value := range c2.data {
		if _, ok := removed[key]; ok {
			continue
		}
		if _, ok := keys[key]; ok {
			continue
		}
		if !process(value) {
			return
		}
	}
}

func (c *Deletable[K, V]) PrettyPrint() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString("Cache Contents:")
	for key, value := range c.data {
		// Handle pointer dereferencing for display purposes
		rv := reflect.ValueOf(value)
		var displayValue interface{} = value
		if rv.Kind() == reflect.Ptr && !rv.IsNil() {
			displayValue = rv.Elem().Interface()
		}

		sb.WriteString(fmt.Sprintf("Key: %+v, Value: %+v,", key, displayValue))
	}

	return sb.String()
}
