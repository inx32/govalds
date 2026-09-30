package concurrency

import (
	"errors"
	"iter"
	"maps"
	"sync"
)

var ErrMapNotFound = errors.New("no such element in Map")

type Map[T comparable, E any] struct {
	m  map[T]E
	mu sync.RWMutex
}

func (m *Map[T, E]) Set(k T, v E) {
	m.mu.Lock()
	m.m[k] = v
	m.mu.Unlock()
}

func (m *Map[T, E]) Update(k T, fn func(v E, exist bool) (new E, keep bool)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, found := m.m[k]
	new, keep := fn(v, found)
	if !keep {
		delete(m.m, k)
		return
	}
	m.m[k] = new
}

func (m *Map[T, E]) Get(k T) (v E, found bool) {
	m.mu.RLock()
	v, found = m.m[k]
	m.mu.RUnlock()
	return
}

func (m *Map[T, E]) MustGet(k T) E {
	v, found := m.Get(k)
	if !found {
		panic(ErrMapNotFound)
	}
	return v
}

func (m *Map[T, E]) DefaultGet(k T, def E) E {
	v, found := m.Get(k)
	if !found {
		return def
	}
	return v
}

func (m *Map[T, E]) GetOrSet(k T, v E) (actual E, found bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if actual, found = m.m[k]; found {
		return
	}
	m.m[k] = v
	return v, false
}

func (m *Map[T, E]) Has(k T) bool {
	m.mu.RLock()
	_, ok := m.m[k]
	m.mu.RUnlock()
	return ok
}

func (m *Map[T, E]) Delete(k T) {
	m.mu.Lock()
	delete(m.m, k)
	m.mu.Unlock()
}

func (m *Map[T, E]) Pop(k T) (v E, found bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, found = m.m[k]
	if found {
		delete(m.m, k)
	}
	return v, found
}

func (m *Map[T, E]) Swap(k T, v E) (old E, existed bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	old, existed = m.m[k]
	m.m[k] = v
	return
}

func (m *Map[T, E]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.m)
}

func (m *Map[T, E]) Clear() {
	m.mu.Lock()
	clear(m.m)
	m.mu.Unlock()
}

func (m *Map[T, E]) IterNoAlloc() iter.Seq2[T, E] {
	return func(yield func(k T, v E) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		for k, v := range m.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (m *Map[T, E]) Iter() iter.Seq2[T, E] {
	m.mu.RLock()
	m2 := maps.Clone(m.m)
	m.mu.RUnlock()

	return func(yield func(k T, v E) bool) {
		for k, v := range m2 {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (m *Map[T, E]) Keys() []T {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]T, len(m.m))
	index := 0

	for k := range m.m {
		keys[index] = k
		index++
	}

	return keys
}

func (m *Map[T, E]) Values() []E {
	m.mu.RLock()
	defer m.mu.RUnlock()

	values := make([]E, len(m.m))
	index := 0

	for _, v := range m.m {
		values[index] = v
		index++
	}

	return values
}

func (m *Map[T, E]) Clone() *Map[T, E] {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return &Map[T, E]{m: maps.Clone(m.m)}
}

func NewMap[T comparable, E any]() *Map[T, E] {
	return &Map[T, E]{m: make(map[T]E)}
}
