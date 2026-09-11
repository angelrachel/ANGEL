package policy

import (
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
)

func TestEngineAllowsOnlyScopedActionAndConsumesBudget(t *testing.T) {
	now := time.Now().UTC()
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-1", Organization: "org-1", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"surface-map"}}}, 1); err != nil {
		t.Fatal(err)
	}
	if decision := engine.Authorize("eng-1", "fixture://lab/web", "surface-map", now); !decision.Allowed {
		t.Fatalf("expected allow: %#v", decision)
	}
	if decision := engine.Authorize("eng-1", "fixture://lab/web", "surface-map", now); decision.Allowed {
		t.Fatal("expected budget denial")
	}
}

func TestEngineDeniesAmbiguousTarget(t *testing.T) {
	now := time.Now().UTC()
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-2", Organization: "org-1", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"surface-map"}}}, 2); err != nil {
		t.Fatal(err)
	}
	if decision := engine.Authorize("eng-2", "https://outside.example", "surface-map", now); decision.Allowed {
		t.Fatal("out-of-scope target allowed")
	}
}
