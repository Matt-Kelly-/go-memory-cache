package store

import (
	"sync"
)

type mutexDecorator struct {
	mutex sync.Mutex // This mutex protects the store
	store Store
}

func WithMutex(store Store) Store {
	return &mutexDecorator{
		store: store,
	}
}

func (s *mutexDecorator) Has(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.store.Has(key)
}

func (s *mutexDecorator) Get(key string) (string, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.store.Get(key)
}

func (s *mutexDecorator) Put(key, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store.Put(key, value)
}

func (s *mutexDecorator) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store.Delete(key)
}
