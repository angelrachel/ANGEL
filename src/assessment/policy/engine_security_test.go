package policy

import (
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
)

func TestEnginePermanentlyDeniesDangerousCapabilities(t *testing.T) {
	now := time.Now().UTC()
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-danger", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"arbitrary-command"}}}, 4); err != nil {
		t.Fatal(err)
	}
	decision := engine.Authorize("eng-danger", "fixture://lab/web", "arbitrary-command", now)
	if decision.Allowed {
		t.Fatal("dangerous capability allowed")
	}
}
func TestEngineMatchesURLScopeBySchemeAndHostname(t *testing.T) {
	now := time.Now().UTC()
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-url", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "url", Value: "https://example.test/base", Actions: []string{"surface-map"}}}, 2); err != nil {
		t.Fatal(err)
	}
	if decision := engine.Authorize("eng-url", "https://example.test/other", "surface-map", now); !decision.Allowed {
		t.Fatalf("same host denied: %#v", decision)
	}
}
