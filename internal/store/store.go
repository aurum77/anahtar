package store

import (
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
}

func (e *entry) isExpired() bool {
	return !e.expiresAt.IsZero() && time.Now().After(e.expiresAt)
}

type Store struct {
	rwmutex sync.RWMutex
	storage map[string]entry
}

func New() *Store {
	storage := make(map[string]entry)
	store := &Store{storage: storage}
	go store.sweep()
	return store
}

func (s *Store) Get(key string) (string, bool) {
	s.rwmutex.RLock()

	entry, ok := s.storage[key]

	if ok && entry.isExpired() {
		defer s.rwmutex.Unlock()

		s.rwmutex.RUnlock() // release read lock
		// between the unlock and the lock , a goroutine might be able
		// to slip in and set the same entry with a later expiry date
		s.rwmutex.Lock() // and acquire full lock for deletion

		entry, ok = s.storage[key]

		if ok && entry.isExpired() {
			delete(s.storage, key)
		} else {
			return entry.value, ok
		}

		return "", false
	}

	s.rwmutex.RUnlock()

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

func (s *Store) sweep() {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		s.removeExpired()
	}
}

func (s *Store) removeExpired() {
	defer s.rwmutex.Unlock()

	s.rwmutex.Lock()

	for key, value := range s.storage {
		if value.isExpired() {
			delete(s.storage, key)
		}
	}
}
