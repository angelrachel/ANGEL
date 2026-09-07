package orchestrator

import (
"context"
"time"
)

type AIDecision struct {
Decision string
Reason   string
}

func NewAIDecision() *AIDecision {
return &AIDecision{}
}

func (a *AIDecision) Analyze(ctx context.Context, input string) {
a.Decision = "Attack"
a.Reason = "Target is vulnerable and reachable."
}

func (a *AIDecision) Decide(ctx context.Context) string {
return a.Decision
}

func (a *AIDecision) GetReason() string {
return a.Reason
}

func (a *AIDecision) Update(ctx context.Context, newDecision string) {
a.Decision = newDecision
}

func (a *AIDecision) GetLastSeen() time.Time {
return time.Now()
}
