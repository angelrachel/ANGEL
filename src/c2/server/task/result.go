package task

import (
	"sync"
	"time"
)

type Result struct {
	TaskID   string
	AgentID  string
	Output   map[string]interface{}
	Status   string
	Finished time.Time
}

type ResultStore struct {
	mu      sync.RWMutex
	results map[string]*Result
}

func NewResultStore() *ResultStore {
	return &ResultStore{results: make(map[string]*Result)}
}

func cloneResult(result *Result) *Result {
	if result == nil {
		return nil
	}
	clone := *result
	if result.Output != nil {
		clone.Output = make(map[string]interface{}, len(result.Output))
		for key, value := range result.Output {
			clone.Output[key] = value
		}
	}
	return &clone
}

func (r *ResultStore) Add(result *Result) {
	if result == nil || result.TaskID == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.results[result.TaskID] = cloneResult(result)
}

func (r *ResultStore) Get(taskID string) *Result {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneResult(r.results[taskID])
}

func (r *ResultStore) List() []*Result {
	r.mu.RLock()
	defer r.mu.RUnlock()
	results := make([]*Result, 0, len(r.results))
	for _, result := range r.results {
		results = append(results, cloneResult(result))
	}
	return results
}
