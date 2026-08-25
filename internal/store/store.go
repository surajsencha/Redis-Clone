package store

import (
	"sync"
	"time"
)

type Store struct {
	data map[string]*Entry
	mu   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]*Entry),
	}
}

func (s *Store) StartActiveExpiry() {
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)

			s.mu.Lock()

			checked := 0

			for k, v := range s.data {
				if v.IsExpired() {
					delete(s.data, k)
				}

				checked++

				if checked >= 20 {
					break
				}
			}

			s.mu.Unlock()
		}
	}()
}

func (s *Store) Set(key string, value interface{}, expireAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = &Entry{
		Value:    value,
		ExpireAt: expireAt,
	}
}

func (s *Store) Get(key string) (*Entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data[key]

	if !ok {
		return nil, false
	}
	if entry.IsExpired() {
		return nil, false
	}
	return entry, true
}

func (s *Store) Del(keys ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	deletedCount := 0
	for _, key := range keys {
		if _, ok := s.data[key]; ok {
			delete(s.data, key)
			deletedCount++
		}
	}
	return deletedCount
}
