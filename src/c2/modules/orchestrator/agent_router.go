package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"

	core "ANGEL/src/c2/orchestrator"
)

type AgentType string

const (
	AgentRecon   AgentType = "recon"
	AgentSQL     AgentType = "sql"
	AgentNoSQL   AgentType = "nosql"
	AgentPost    AgentType = "post"
	AgentLateral AgentType = "lateral"
	AgentDestroy AgentType = "destroy"
)

type AgentTask struct {
	Type        AgentType
	Target      string
	Command     string
	Technique   string
	Class       core.ActionClass
	ApprovalID  string
	RequestedBy string
}

type AgentResult struct {
	Type       AgentType
	Success    bool
	Output     string
	Authorized bool
	Reason     string
}

type AgentRouter struct {
	Mu      sync.Mutex
	Results chan AgentResult
	Scope   *core.EngagementScope
}

func NewAgentRouter() *AgentRouter {
	return &AgentRouter{Results: make(chan AgentResult, 100)}
}

func (a *AgentRouter) SetScope(scope core.EngagementScope) {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	a.Scope = &scope
}

func (a *AgentRouter) Route(ctx context.Context, task AgentTask) (AgentResult, error) {
	if err := ctx.Err(); err != nil {
		return AgentResult{Type: task.Type, Reason: err.Error()}, err
	}

	a.Mu.Lock()
	scope := a.Scope
	a.Mu.Unlock()
	if scope == nil {
		return AgentResult{Type: task.Type, Reason: "engagement scope is not configured"}, fmt.Errorf("engagement scope is not configured")
	}

	decision := scope.Authorize(core.ActionRequest{
		Target:      task.Target,
		Technique:   task.Technique,
		Class:       task.Class,
		ApprovalID:  task.ApprovalID,
		RequestedBy: task.RequestedBy,
	}, time.Now().UTC())
	if !decision.Allowed {
		result := AgentResult{Type: task.Type, Authorized: false, Reason: decision.Reason}
		return result, fmt.Errorf("authorization denied: %s", decision.Reason)
	}

	if !knownAgent(task.Type) {
		return AgentResult{Type: task.Type, Authorized: true, Reason: "unknown agent type"}, fmt.Errorf("unknown agent type")
	}

	// This router deliberately acknowledges a policy-approved simulation route;
	// it does not execute commands, create implants, or perform target actions.
	result := AgentResult{
		Type:       task.Type,
		Success:    true,
		Authorized: true,
		Output:     "simulation route accepted; no target action executed",
		Reason:     decision.Reason,
	}
	select {
	case a.Results <- result:
	default:
	}
	return result, nil
}

func knownAgent(agent AgentType) bool {
	switch agent {
	case AgentRecon, AgentSQL, AgentNoSQL, AgentPost, AgentLateral, AgentDestroy:
		return true
	default:
		return false
	}
}

func (a *AgentRouter) GetLastSeen() time.Time {
	return time.Now().UTC()
}
