package store

import (
	"sync"
)

type rwMutexDecorator struct {
	mutex sync.RWMutex // This mutex protects the store
	store Store
}

func WithRWMutex(store Store) Store {
	return &rwMutexDecorator{
		store: store,
	}
}

func (s *rwMutexDecorator) Has(key string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.store.Has(key)
}

func (s *rwMutexDecorator) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.store.Get(key)
}

func (s *rwMutexDecorator) Put(key, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store.Put(key, value)
}

func (s *rwMutexDecorator) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store.Delete(key)
}
