package orchestrator

import (
	"context"
	"strings"
	"testing"
	"time"

	core "ANGEL/src/assessment/orchestrator"
)

func TestAssessmentRouterAcceptsScopedWorker(t *testing.T) {
	now := time.Now().UTC()
	router := NewAgentRouter()
	router.SetScope(core.EngagementScope{ID: "eng", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute), Targets: []string{"fixture://lab/web"}, AllowedTechniques: []string{"surface-map"}, RulesOfEngagement: "roe"})
	result, err := router.Route(context.Background(), AgentTask{Type: WebAssessment, Target: "fixture://lab/web", Technique: "surface-map", Class: core.ActionSimulation, RequestedBy: "operator"})
	if err != nil || !result.Success || !result.Authorized || !result.Simulation {
		t.Fatalf("result=%#v err=%v", result, err)
	}
	if !strings.Contains(result.Output, "no arbitrary target action") {
		t.Fatalf("output=%q", result.Output)
	}
}

func TestAssessmentRouterRejectsUnknownWorker(t *testing.T) {
	now := time.Now().UTC()
	router := NewAgentRouter()
	router.SetScope(core.EngagementScope{ID: "eng", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute), Targets: []string{"fixture://lab/web"}, AllowedTechniques: []string{"surface-map"}, RulesOfEngagement: "roe"})
	if _, err := router.Route(context.Background(), AgentTask{Type: AgentType("arbitrary"), Target: "fixture://lab/web", Technique: "surface-map", Class: core.ActionSimulation}); err == nil {
		t.Fatal("unknown worker accepted")
	}
}
