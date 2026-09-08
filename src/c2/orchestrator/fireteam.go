package orchestrator

import (
"context"
"sync"
)

type Fireteam struct {
mu      sync.Mutex
tasks   []Task
results []Result
}

func NewFireteam() *Fireteam {
return &Fireteam{}
}

func (f *Fireteam) Launch(ctx context.Context, router *Router, tasks []Task) []Result {
f.mu.Lock()
defer f.mu.Unlock()
f.tasks = tasks
f.results = []Result{}
for _, task := range tasks {
select {
case <-ctx.Done():
return f.results
default:
result := router.Dispatch(task)
f.results = append(f.results, result)
}
}
return f.results
}

func (f *Fireteam) GetResults() []Result {
f.mu.Lock()
defer f.mu.Unlock()
return f.results
}
