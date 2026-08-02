package store

import "sync/atomic"

type PointerKVStore[K comparable, V any] struct {
	value atomic.Pointer[map[K]V]
}

func NewPointerKVStore[K comparable, V any]() *PointerKVStore[K, V] {
	s := &PointerKVStore[K, V]{}
	empty := make(map[K]V)
	s.value.Store(&empty)

	return s
}

func (s *PointerKVStore[K, V]) ReadByKey(key K) V {
	value := *s.value.Load()

	return value[key]
}

func (s *PointerKVStore[K, V]) Set(value map[K]V) {
	s.value.Store(&value)
}
