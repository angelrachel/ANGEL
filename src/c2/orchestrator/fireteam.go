package orchestrator

import (
"context"
"sync"
)

type Fireteam struct {
Agents     []*AgentRouter
TaskQueue  chan AgentTask
ResultChan chan AgentResult
}

func NewFireteam(agents []*AgentRouter) *Fireteam {
return &Fireteam{
Agents:     agents,
TaskQueue:  make(chan AgentTask, 100),
ResultChan: make(chan AgentResult, 100),
}
}

func (f *Fireteam) AssignTask(task AgentTask) {
f.TaskQueue <- task
}

func (f *Fireteam) Execute(ctx context.Context) {
var wg sync.WaitGroup
for _, agent := range f.Agents {
wg.Add(1)
go func(a *AgentRouter) {
defer wg.Done()
for task := range f.TaskQueue {
result, _ := a.Route(ctx, task)
f.ResultChan <- result
}
}(agent)
}
wg.Wait()
close(f.ResultChan)
}

func (f *Fireteam) GetResults() chan AgentResult {
return f.ResultChan
}

func (f *Fireteam) Stop() {
close(f.TaskQueue)
}
