package policy

import (
	"ANGEL/src/assessment/domain"
	"testing"
	"time"
)

func TestHostnameWildcardRequiresLabelBoundary(t *testing.T) {
	now := time.Now().UTC()
	engine := NewEngine()
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-host", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "hostname", Value: "*.example.com", Actions: []string{"surface-map"}}}, 3); err != nil {
		t.Fatal(err)
	}
	if !engine.Authorize("eng-host", "api.example.com", "surface-map", now).Allowed {
		t.Fatal("subdomain was rejected")
	}
	if engine.Authorize("eng-host", "badexample.com", "surface-map", now).Allowed {
		t.Fatal("suffix collision was accepted")
	}
}

func TestURLPathAndPortScope(t *testing.T) {
	now := time.Now().UTC()
	engine := NewEngine()
	scope := []domain.ScopeEntry{{Kind: "url", Value: "https://example.com/api", Ports: []int{443}, Actions: []string{"api-contract"}}, {Kind: "path", Value: "/api", Actions: []string{"path-check"}}}
	if err := engine.RegisterEngagement(domain.Engagement{ID: "eng-url", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, scope, 3); err != nil {
		t.Fatal(err)
	}
	if !engine.Authorize("eng-url", "https://example.com/api/v1", "api-contract", now).Allowed {
		t.Fatal("scoped path was rejected")
	}
	if engine.Authorize("eng-url", "http://example.com/api", "api-contract", now).Allowed {
		t.Fatal("wrong scheme/default port was accepted")
	}
	if !engine.Authorize("eng-url", "https://example.com/api/v1", "path-check", now).Allowed {
		t.Fatal("explicit path scope was rejected")
	}
}
