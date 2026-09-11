package policy

import (
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
)

// Test semua capability yang ditolak tetap ditolak
func TestEnginePermanentlyDeniesAllDangerousCapabilities(t *testing.T) {
	deniedCapabilities := []string{
		"credential-collection",
		"credential-extraction",
		"persistence",
		"destructive-write",
		"log-deletion",
		"covert-channel",
		"process-injection",
		"evasion",
		"data-exfiltration",
		"arbitrary-command",
	}

	now := time.Now().UTC()
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-cap-test", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: deniedCapabilities}}, 10); err != nil {
		t.Fatal(err)
	}

	for _, cap := range deniedCapabilities {
		decision := engine.Authorize("eng-cap-test", "fixture://lab/web", cap, now)
		if decision.Allowed {
			t.Fatalf("dangerous capability %q should be denied but was allowed", cap)
		}
		if decision.Reason == "" {
			t.Fatalf("dangerous capability %q should have a denial reason", cap)
		}
	}
}

// Test bahwa capability yang diizinkan tetap diizinkan bila dalam scope
func TestEngineAllowsSafeCapabilitiesWithinScope(t *testing.T) {
	safeCapabilities := []string{
		"surface-map",
		"tls-assessment",
		"api-contract",
		"authorization-matrix",
		"evidence-collection",
		"dependency-inventory",
		"detection-validation",
		"synthetic-canary",
		"lab-proof",
	}

	now := time.Now().UTC()
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-safe-cap", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: safeCapabilities}}, 10); err != nil {
		t.Fatal(err)
	}

	for _, cap := range safeCapabilities {
		decision := engine.Authorize("eng-safe-cap", "fixture://lab/web", cap, now)
		if !decision.Allowed {
			t.Fatalf("safe capability %q should be allowed within scope but was denied: %#v", cap, decision)
		}
	}
}

// Test bahwa scope deny (target di luar scope) tetap ditolak
func TestEngineDeniesTargetOutsideScope(t *testing.T) {
	now := time.Now().UTC()
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-scope-test", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web-01", Actions: []string{"surface-map"}}}, 5); err != nil {
		t.Fatal(err)
	}

	outOfScopeTargets := []string{
		"fixture://lab/web-02",
		"fixture://other/target",
		"https://example.com",
		"https://outside.invalid",
	}

	for _, target := range outOfScopeTargets {
		decision := engine.Authorize("eng-scope-test", target, "surface-map", now)
		if decision.Allowed {
			t.Fatalf("target %q should be denied (out of scope) but was allowed", target)
		}
	}
}

// Test bahwa action di luar scope actions ditolak
func TestEngineDeniesActionOutsideScopeActions(t *testing.T) {
	now := time.Now().UTC()
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-action-test", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"surface-map"}}}, 5); err != nil {
		t.Fatal(err)
	}

	outOfScopeActions := []string{
		"api-contract",
		"tls-assessment",
		"arbitrary-command",
	}

	for _, action := range outOfScopeActions {
		decision := engine.Authorize("eng-action-test", "fixture://lab/web", action, now)
		if decision.Allowed {
			t.Fatalf("action %q should be denied (not in scope actions) but was allowed", action)
		}
	}
}

// Test budget exhaustion
func TestEngineDeniesAfterBudgetExhausted(t *testing.T) {
	now := time.Now().UTC()
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-budget-test", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"surface-map"}}}, 2); err != nil {
		t.Fatal(err)
	}

	// Penggunaan budget 2
	if decision := engine.Authorize("eng-budget-test", "fixture://lab/web", "surface-map", now); !decision.Allowed {
		t.Fatalf("first usage should be allowed, got: %#v", decision)
	}
	if decision := engine.Authorize("eng-budget-test", "fixture://lab/web", "surface-map", now); !decision.Allowed {
		t.Fatalf("second usage should be allowed, got: %#v", decision)
	}
	// Budget habis, ditolak
	if decision := engine.Authorize("eng-budget-test", "fixture://lab/web", "surface-map", now); decision.Allowed {
		t.Fatal("third usage should be denied (budget exhausted)")
	}
}

// Test engagement expired - use valid engagement but request with past time
func TestEngineDeniesExpiredEngagement(t *testing.T) {
	now := time.Now().UTC()
	past := now.Add(-time.Hour)
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-expired", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Hour)}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"surface-map"}}}, 5); err != nil {
		t.Fatal(err)
	}

	// Minta authorize dengan waktu di masa lalu yang sudah melewati ends_at
	decision := engine.Authorize("eng-expired", "fixture://lab/web", "surface-map", past)
	if decision.Allowed {
		t.Fatal("expired time should be denied")
	}
}

// Test engagement emergency stop
func TestEngineDeniesWhenEngagementStopped(t *testing.T) {
	now := time.Now().UTC()
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-stopped", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Hour), Stopped: true}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"surface-map"}}}, 5); err != nil {
		t.Fatal(err)
	}

	decision := engine.Authorize("eng-stopped", "fixture://lab/web", "surface-map", now)
	if decision.Allowed {
		t.Fatal("stopped engagement should be denied")
	}
}

// Test fixture target enforcement
func TestEngineRequiresFixtureTarget(t *testing.T) {
	now := time.Now().UTC()
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-fixture-test", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"surface-map"}}}, 5); err != nil {
		t.Fatal(err)
	}

	nonFixtureTargets := []string{
		"https://example.com",
		"http://internal.corp",
		"192.168.1.1",
		"example.com",
	}

	for _, target := range nonFixtureTargets {
		decision := engine.Authorize("eng-fixture-test", target, "surface-map", now)
		if decision.Allowed {
			t.Fatalf("non-fixture target %q should be denied", target)
		}
	}
}
