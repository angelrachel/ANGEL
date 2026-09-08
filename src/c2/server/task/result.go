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
mu      sync.Mutex
results map[string]*Result
}

func NewResultStore() *ResultStore {
return &ResultStore{results: make(map[string]*Result)}
}

func (r *ResultStore) Add(result *Result) {
r.mu.Lock()
defer r.mu.Unlock()
r.results[result.TaskID] = result
}

func (r *ResultStore) Get(taskID string) *Result {
r.mu.Lock()
defer r.mu.Unlock()
return r.results[taskID]
}

func (r *ResultStore) List() []*Result {
r.mu.Lock()
defer r.mu.Unlock()
var results []*Result
for _, result := range r.results {
results = append(results, result)
}
return results
}
