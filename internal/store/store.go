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

func (s *Store) RPush(key string, values ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]

	// Key doesn't exist or has expired
	if !ok || entry.IsExpired() {
		list := []string{}

		list = append(list, values...)

		s.data[key] = &Entry{
			Value: list,
		}

		return len(list)
	}

	// Key exists, so its value should be a list
	list, ok := entry.Value.([]string)

	if !ok {
		panic("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	// Add new values to the existing list
	list = append(list, values...)

	// Save updated list
	entry.Value = list
	s.data[key] = entry

	return len(list)
}

func (s *Store) LRange(key string, start, stop int) ([]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data[key]

	if !ok || entry.IsExpired() {
		return nil, false
	}

	list, ok := entry.Value.([]string)

	if !ok {
		return nil, false
	}
	n := len(list)

	// Normalize negative indexes
	if start < 0 {
		start = n + start
	}

	if stop < 0 {
		stop = n + stop
	}

	// Clamp to bounds
	if start < 0 {
		start = 0
	}
	// Clamp to bounds (stop should never be higher than the last index)
	if stop >= n {
		stop = n - 1
	}

	// Invalid range
	if start > stop {
		return []string{}, true
	}

	// Add 1 to stop because Go slices are exclusive!
	return list[start : stop+1], true

}
