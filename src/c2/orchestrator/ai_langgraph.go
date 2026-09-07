package orchestrator

import (
"context"
"time"
)

type AILangGraph struct {
State    map[string]interface{}
Current  string
Start    string
End      string
}

func NewAILangGraph(start, end string) *AILangGraph {
return &AILangGraph{
State:   make(map[string]interface{}),
Start:   start,
End:     end,
Current: start,
}
}

func (a *AILangGraph) SetState(key string, value interface{}) {
a.State[key] = value
}

func (a *AILangGraph) GetState(key string) interface{} {
return a.State[key]
}

func (a *AILangGraph) Next(ctx context.Context) string {
return a.End
}

func (a *AILangGraph) Traverse(ctx context.Context) []string {
var path []string
current := a.Start
for current != "" {
path = append(path, current)
if current == a.End {
break
}
current = a.Next(ctx)
}
return path
}

func (a *AILangGraph) GetCurrent() string {
return a.Current
}

func (a *AILangGraph) GetLastSeen() time.Time {
return time.Now()
}
