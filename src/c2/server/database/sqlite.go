package database

import (
	"sync"
	"time"
)

type Store struct {
	mu         sync.RWMutex
	data       map[string]interface{}
	generation map[string]uint64
}

func NewStore() *Store {
	return &Store{
		data:       make(map[string]interface{}),
		generation: make(map[string]uint64),
	}
}

func (s *Store) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.generation[key]++
	s.data[key] = value
}

func (s *Store) Get(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[key]
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.generation[key]++
	delete(s.data, key)
}

// List returns a snapshot so callers cannot mutate store state without the lock.
func (s *Store) List() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snapshot := make(map[string]interface{}, len(s.data))
	for key, value := range s.data {
		snapshot[key] = value
	}
	return snapshot
}

func (s *Store) SetWithExpiry(key string, value interface{}, expiry time.Duration) {
	s.mu.Lock()
	s.generation[key]++
	generation := s.generation[key]
	s.data[key] = value
	s.mu.Unlock()

	if expiry <= 0 {
		return
	}
	go func() {
		timer := time.NewTimer(expiry)
		defer timer.Stop()
		<-timer.C
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.generation[key] == generation {
			delete(s.data, key)
			s.generation[key]++
		}
	}()
}
