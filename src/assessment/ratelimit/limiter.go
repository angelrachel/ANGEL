package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

type bucket struct {
	tokens  float64
	updated time.Time
}
type Limiter struct {
	mu      sync.Mutex
	rate    float64
	burst   float64
	buckets map[string]bucket
}

func New(ratePerSecond float64, burst int) *Limiter {
	if ratePerSecond <= 0 {
		ratePerSecond = 1
	}
	if burst < 1 {
		burst = 1
	}
	return &Limiter{rate: ratePerSecond, burst: float64(burst), buckets: map[string]bucket{}}
}
func (l *Limiter) Allow(key string, now time.Time) error {
	if key == "" {
		return fmt.Errorf("rate limit key is required")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	state, ok := l.buckets[key]
	if !ok {
		state = bucket{tokens: l.burst, updated: now}
	}
	elapsed := now.Sub(state.updated).Seconds()
	if elapsed > 0 {
		state.tokens += elapsed * l.rate
		if state.tokens > l.burst {
			state.tokens = l.burst
		}
		state.updated = now
	}
	if state.tokens < 1 {
		l.buckets[key] = state
		return fmt.Errorf("rate limit exceeded")
	}
	state.tokens--
	l.buckets[key] = state
	return nil
}
func (l *Limiter) Reset(key string) { l.mu.Lock(); delete(l.buckets, key); l.mu.Unlock() }
