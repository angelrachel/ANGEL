package agents

import (
"context"
"time"
)

type TaskStatus int

const (
StatusPending TaskStatus = iota
StatusRunning
StatusCompleted
StatusFailed
)

type Task struct {
ID        string
Type      string
Payload   map[string]interface{}
CreatedAt time.Time
}

type Result struct {
TaskID    string
Status    TaskStatus
Data      map[string]interface{}
Error     string
StartedAt time.Time
EndedAt   time.Time
}

type Agent interface {
Execute(ctx context.Context, task Task) (Result, error)
GetName() string
}
