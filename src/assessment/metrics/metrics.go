package metrics

import (
	"sync/atomic"
	"time"
)

type Snapshot struct {
	Requests       uint64 `json:"requests"`
	Errors         uint64 `json:"errors"`
	Active         int64  `json:"active"`
	TotalLatencyMS uint64 `json:"total_latency_ms"`
}
type Collector struct {
	requests atomic.Uint64
	errors   atomic.Uint64
	active   atomic.Int64
	latency  atomic.Uint64
}

func New() *Collector       { return &Collector{} }
func (c *Collector) Start() { c.requests.Add(1); c.active.Add(1) }
func (c *Collector) Finish(status int, started time.Time) {
	c.active.Add(-1)
	if status >= 400 {
		c.errors.Add(1)
	}
	c.latency.Add(uint64(time.Since(started).Milliseconds()))
}
func (c *Collector) Snapshot() Snapshot {
	return Snapshot{Requests: c.requests.Load(), Errors: c.errors.Load(), Active: c.active.Load(), TotalLatencyMS: c.latency.Load()}
}
