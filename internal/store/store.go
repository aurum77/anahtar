package store

import (
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
}

type Store struct {
	rwmutex sync.RWMutex
	storage map[string]entry
}

func New() *Store {
	storage := make(map[string]entry)
	return &Store{storage: storage}
}

func (s *Store) Get(key string) (string, bool) {
	defer s.rwmutex.RUnlock()

	s.rwmutex.RLock()

	entry, ok := s.storage[key]

	return entry.value, ok
}

func (s *Store) Set(key string, value string, expiresAt time.Time) {
	defer s.rwmutex.Unlock()

	s.rwmutex.Lock()

	entry := entry{value, expiresAt}

	s.storage[key] = entry
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
