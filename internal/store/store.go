package store

import (
	"sync"
)

type Store struct {
	rwmutex sync.RWMutex
	storage map[string]string
}

func New() *Store {
	storage := make(map[string]string)
	return &Store{storage: storage}
}

func (s *Store) Get(key string) (string, bool) {
	defer s.rwmutex.RUnlock()

	s.rwmutex.RLock()

	value, ok := s.storage[key]

	return value, ok
}

func (s *Store) Set(key string, value string) {
	defer s.rwmutex.Unlock()

	s.rwmutex.Lock()

	s.storage[key] = value
}

func (s *Store) Del(key string) bool {
	defer s.rwmutex.Unlock()

	s.rwmutex.Lock()

	_, ok := s.storage[key]

	if !ok {
		return false
	}

	delete(s.storage, key)

	return true
}
