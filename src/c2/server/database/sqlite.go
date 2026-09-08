package database

import (
"sync"
"time"
)

type Store struct {
mu   sync.Mutex
data map[string]interface{}
}

func NewStore() *Store {
return &Store{data: make(map[string]interface{})}
}

func (s *Store) Set(key string, value interface{}) {
s.mu.Lock()
defer s.mu.Unlock()
s.data[key] = value
}

func (s *Store) Get(key string) interface{} {
s.mu.Lock()
defer s.mu.Unlock()
return s.data[key]
}

func (s *Store) Delete(key string) {
s.mu.Lock()
defer s.mu.Unlock()
delete(s.data, key)
}

func (s *Store) List() map[string]interface{} {
s.mu.Lock()
defer s.mu.Unlock()
return s.data
}

func (s *Store) SetWithExpiry(key string, value interface{}, expiry time.Duration) {
s.Set(key, value)
go func() {
time.Sleep(expiry)
s.Delete(key)
}()
}
