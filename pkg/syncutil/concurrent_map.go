package syncutil

import "sync"

type ConcurrentMap[K comparable, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Remove(key K)
}

type concurrentMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewConcurrentMap[K comparable, V any]() ConcurrentMap[K, V] {
	return &concurrentMap[K, V]{
		data: make(map[K]V),
	}
}

func (m *concurrentMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *concurrentMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.data[key]
	return value, ok
}

func (m *concurrentMap[K, V]) Remove(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}
