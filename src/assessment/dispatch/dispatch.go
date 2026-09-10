package dispatch

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Event struct {
	Type    string            `json:"type"`
	JobID   string            `json:"job_id"`
	At      time.Time         `json:"at"`
	Payload map[string]string `json:"payload"`
}
type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]chan Event
}

func NewBus() *Bus { return &Bus{subscribers: map[string][]chan Event{}} }
func (b *Bus) Subscribe(eventType string, buffer int) (<-chan Event, func()) {
	if buffer < 1 {
		buffer = 1
	}
	ch := make(chan Event, buffer)
	b.mu.Lock()
	b.subscribers[eventType] = append(b.subscribers[eventType], ch)
	b.mu.Unlock()
	return ch, func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		items := b.subscribers[eventType]
		for i, item := range items {
			if item == ch {
				b.subscribers[eventType] = append(items[:i], items[i+1:]...)
				close(ch)
				break
			}
		}
	}
}
func (b *Bus) Publish(event Event) {
	b.mu.RLock()
	subs := append([]chan Event(nil), b.subscribers[event.Type]...)
	b.mu.RUnlock()
	for _, ch := range subs {
		select {
		case ch <- event:
		default:
		}
	}
}

type Work func(context.Context) error
type Pool struct {
	queue   chan Work
	workers int
}

func NewPool(workers, capacity int) *Pool {
	if workers < 1 {
		workers = 1
	}
	if capacity < workers {
		capacity = workers
	}
	return &Pool{queue: make(chan Work, capacity), workers: workers}
}
func (p *Pool) Start(ctx context.Context) *sync.WaitGroup {
	group := &sync.WaitGroup{}
	for i := 0; i < p.workers; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case work := <-p.queue:
					if work != nil {
						_ = work(ctx)
					}
				}
			}
		}()
	}
	return group
}
func (p *Pool) Submit(ctx context.Context, work Work) error {
	if work == nil {
		return fmt.Errorf("work is required")
	}
	select {
	case p.queue <- work:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (p *Pool) Capacity() int { return cap(p.queue) }
