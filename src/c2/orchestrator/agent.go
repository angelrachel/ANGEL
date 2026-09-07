package orchestrator

import (
"context"
"time"
)

type Agent struct {
ID       string
Type     AgentType
Task     AgentTask
Result   AgentResult
LastSeen time.Time
}

func NewAgent(id string, agentType AgentType) *Agent {
return &Agent{
ID:       id,
Type:     agentType,
LastSeen: time.Now(),
}
}

func (a *Agent) Execute(ctx context.Context) AgentResult {
a.LastSeen = time.Now()
return AgentResult{
Type:    a.Type,
Target:  a.Task.Target,
Payload: a.Task.Command,
Success: true,
Time:    time.Now(),
}
}

func (a *Agent) UpdateLastSeen() {
a.LastSeen = time.Now()
}

func (a *Agent) GetID() string {
return a.ID
}

func (a *Agent) GetType() AgentType {
return a.Type
}

func (a *Agent) GetLastSeen() time.Time {
return a.LastSeen
}
