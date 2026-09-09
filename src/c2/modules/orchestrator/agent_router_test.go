package orchestrator

import (
	"context"
	"strings"
	"testing"
	"time"

	core "ANGEL/src/c2/orchestrator"
)

func routerScope(now time.Time) core.EngagementScope {
	return core.EngagementScope{
		ID:                "eng-router-001",
		Authorized:        true,
		StartsAt:          now.Add(-time.Minute),
		EndsAt:            now.Add(time.Minute),
		Targets:           []string{"lab.example"},
		AllowedTechniques: []string{"simulated-validation"},
		RulesOfEngagement: "roe-router-001",
	}
}

func TestAgentRouterFailsClosedWithoutScope(t *testing.T) {
	router := NewAgentRouter()
	_, err := router.Route(context.Background(), AgentTask{
		Type: AgentRecon, Target: "lab.example", Technique: "simulated-validation",
		Class: core.ActionSimulation, RequestedBy: "operator",
	})
	if err == nil || !strings.Contains(err.Error(), "scope is not configured") {
		t.Fatalf("expected missing scope rejection, got %v", err)
	}
}

func TestAgentRouterAcceptsAuthorizedSimulationOnly(t *testing.T) {
	now := time.Now().UTC()
	router := NewAgentRouter()
	router.SetScope(routerScope(now))

	result, err := router.Route(context.Background(), AgentTask{
		Type: AgentRecon, Target: "LAB.EXAMPLE", Technique: "SIMULATED-VALIDATION",
		Class: core.ActionSimulation, RequestedBy: "operator",
	})
	if err != nil || !result.Success || !result.Authorized {
		t.Fatalf("expected authorized simulation route, result=%#v err=%v", result, err)
	}
	if !strings.Contains(result.Output, "no target action executed") {
		t.Fatalf("expected non-executing route result, got %q", result.Output)
	}
}

func TestAgentRouterRejectsOutOfScopeAndActiveWithoutApproval(t *testing.T) {
	now := time.Now().UTC()
	router := NewAgentRouter()
	router.SetScope(routerScope(now))

	cases := []AgentTask{
		{Type: AgentRecon, Target: "outside.example", Technique: "simulated-validation", Class: core.ActionSimulation, RequestedBy: "operator"},
		{Type: AgentRecon, Target: "lab.example", Technique: "simulated-validation", Class: core.ActionActive, RequestedBy: "operator"},
	}
	for _, task := range cases {
		result, err := router.Route(context.Background(), task)
		if err == nil || result.Authorized {
			t.Fatalf("expected task rejection, result=%#v err=%v", result, err)
		}
	}
}
