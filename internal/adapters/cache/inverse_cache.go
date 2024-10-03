package cache

import (
	"fmt"
	"strings"
	"sync"
)

type InverseCache[K comparable, V comparable] struct {
	mu      sync.Mutex
	direct  map[K]V
	inverse map[V]K
}

func NewInverseCache[K comparable, V comparable]() *InverseCache[K, V] {
	return &InverseCache[K, V]{
		mu:      sync.Mutex{},
		direct:  make(map[K]V),
		inverse: make(map[V]K),
	}
}

func (i *InverseCache[K, V]) Add(a K, b V) {
	i.mu.Lock()
	i.direct[a] = b
	i.inverse[b] = a
	i.mu.Unlock()
}

func (i *InverseCache[K, V]) GetInverse(b V) (K, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	val, ok := i.inverse[b]
	return val, ok
}

func (i *InverseCache[K, V]) GetDirect(a K) (V, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	val, ok := i.direct[a]
	return val, ok
}

func (i *InverseCache[K, V]) DirectContains(a K) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	_, ok := i.direct[a]
	return ok
}

func (i *InverseCache[K, V]) InverseContains(b V) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	_, ok := i.inverse[b]
	return ok
}

func (i *InverseCache[K, V]) DeleteDirect(a K) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if val, ok := i.direct[a]; ok {
		delete(i.inverse, val)
	}
	delete(i.direct, a)
}

func (i *InverseCache[K, V]) DeleteInverse(b V) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if val, ok := i.inverse[b]; ok {
		delete(i.direct, val)
	}
	delete(i.inverse, b)
}

func (i *InverseCache[K, V]) Delete(a K, b V) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if val, ok := i.inverse[b]; ok {
		delete(i.direct, val)
	}
	delete(i.inverse, b)

	if val, ok := i.direct[a]; ok {
		delete(i.inverse, val)
	}
	delete(i.direct, a)
}

func (i *InverseCache[K, V]) PrettyPrint() string {
	i.mu.Lock()
	defer i.mu.Unlock()

	var sb strings.Builder
	sb.WriteString("Direct Cache:")
	for key, value := range i.direct {
		sb.WriteString(fmt.Sprintf("Key: %+v, Value: %+v,", key, value))
	}

	sb.WriteString("Inverse Cache:")
	for key, value := range i.inverse {
		sb.WriteString(fmt.Sprintf("Value: %+v, Key: %+v,", key, value))
	}

	return sb.String()
}
