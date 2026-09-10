package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID          string    `json:"id"`
	Attempts    int       `json:"attempts"`
	MaxAttempts int       `json:"max_attempts"`
	NotBefore   time.Time `json:"not_before"`
}
type Handler func(context.Context, Job) error
type Scheduler struct {
	mu         sync.Mutex
	queue      []Job
	workers    int
	handler    Handler
	retryDelay time.Duration
}

func New(workers int, handler Handler) *Scheduler {
	if workers < 1 {
		workers = 1
	}
	return &Scheduler{workers: workers, handler: handler, retryDelay: 100 * time.Millisecond}
}
func (s *Scheduler) Enqueue(job Job) error {
	if job.ID == "" || s.handler == nil {
		return fmt.Errorf("job and handler are required")
	}
	if job.MaxAttempts < 1 {
		job.MaxAttempts = 1
	}
	s.mu.Lock()
	s.queue = append(s.queue, job)
	s.mu.Unlock()
	return nil
}
func (s *Scheduler) Pending() int { s.mu.Lock(); defer s.mu.Unlock(); return len(s.queue) }
func (s *Scheduler) Run(ctx context.Context) error {
	if s.handler == nil {
		return fmt.Errorf("handler is required")
	}
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		job, ok := s.popReady(time.Now())
		if !ok {
			if s.Pending() == 0 {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(10 * time.Millisecond):
				continue
			}
		}
		err := s.handler(ctx, job)
		if err == nil {
			continue
		}
		job.Attempts++
		if job.Attempts >= job.MaxAttempts {
			continue
		}
		job.NotBefore = time.Now().Add(s.retryDelay * time.Duration(job.Attempts))
		s.mu.Lock()
		s.queue = append(s.queue, job)
		s.mu.Unlock()
	}
}
func (s *Scheduler) popReady(now time.Time) (Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, item := range s.queue {
		if item.NotBefore.IsZero() || !now.Before(item.NotBefore) {
			s.queue = append(s.queue[:i], s.queue[i+1:]...)
			return item, true
		}
	}
	return Job{}, false
}
