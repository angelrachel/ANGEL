package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"

	core "ANGEL/src/assessment/orchestrator"
)

type AgentType string

const (
	AssetDiscovery     AgentType = "asset-discovery"
	NetworkAssessment  AgentType = "network-assessment"
	WebAssessment      AgentType = "web-assessment"
	APIAssessment      AgentType = "api-assessment"
	IdentityAssessment AgentType = "identity-assessment"
	CloudAssessment    AgentType = "cloud-assessment"
	EndpointAssessment AgentType = "endpoint-assessment"
	EvidenceCollection AgentType = "evidence-collection"
)

type AgentTask struct {
	Type        AgentType
	Target      string
	Technique   string
	Class       core.ActionClass
	ApprovalID  string
	RequestedBy string
}
type AgentResult struct {
	Type       AgentType
	Success    bool
	Simulation bool
	Output     string
	Authorized bool
	Reason     string
}
type AgentRouter struct {
	Mu      sync.Mutex
	Results chan AgentResult
	Scope   *core.EngagementScope
}

func NewAgentRouter() *AgentRouter { return &AgentRouter{Results: make(chan AgentResult, 100)} }
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
	if task.Class != core.ActionSimulation {
		return AgentResult{Type: task.Type, Reason: "assessment router accepts observe or simulation actions only"}, fmt.Errorf("assessment router accepts observe or simulation actions only")
	}
	decision := scope.Authorize(core.ActionRequest{Target: task.Target, Technique: task.Technique, Class: task.Class, ApprovalID: task.ApprovalID, RequestedBy: task.RequestedBy}, time.Now().UTC())
	if !decision.Allowed {
		return AgentResult{Type: task.Type, Authorized: false, Reason: decision.Reason}, fmt.Errorf("authorization denied: %s", decision.Reason)
	}
	if !knownAgent(task.Type) {
		return AgentResult{Type: task.Type, Authorized: true, Reason: "unknown assessment worker"}, fmt.Errorf("unknown assessment worker")
	}
	result := AgentResult{Type: task.Type, Success: true, Simulation: true, Authorized: true, Output: "policy-approved assessment route accepted; no arbitrary target action executed", Reason: decision.Reason}
	select {
	case a.Results <- result:
	default:
	}
	return result, nil
}
func knownAgent(agent AgentType) bool {
	switch agent {
	case AssetDiscovery, NetworkAssessment, WebAssessment, APIAssessment, IdentityAssessment, CloudAssessment, EndpointAssessment, EvidenceCollection:
		return true
	default:
		return false
	}
}
func (a *AgentRouter) GetLastSeen() time.Time { return time.Now().UTC() }
