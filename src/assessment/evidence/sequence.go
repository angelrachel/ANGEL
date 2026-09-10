package evidence

import "sync"

type Sequence struct {
	mu sync.Mutex
	n  int64
}

func NewSequence() *Sequence {
	return &Sequence{}
}

func (s *Sequence) Next() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.n++
	return s.n
}
