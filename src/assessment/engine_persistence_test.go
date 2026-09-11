package main

import (
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
	"ANGEL/src/assessment/governance"
)

func TestEngineRestoresDurableSignedJobState(t *testing.T) {
	path := t.TempDir() + "/angel-state.json"
	t.Setenv("ANGEL_STATE_PATH", path)
	now := time.Now().UTC()
	first, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	engagement := domain.Engagement{ID: "persist-eng", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Hour)}
	if err := first.policyEngine.RegisterEngagement(engagement, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"surface-map"}}}, 2); err != nil {
		t.Fatal(err)
	}
	job, err := first.jobController.Create(engagement.ID, "surface-map", "fixture://lab/web", "surface-map", now)
	if err != nil {
		t.Fatal(err)
	}
	first.audit.Append("operator", "PERSISTENCE_TEST", "engagement", engagement.ID, now)
	first.persistState()
	second, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	restored, ok := second.jobController.Get(job.ID)
	if !ok || restored.ID != job.ID || !second.jobController.Verify(restored) {
		t.Fatalf("restored=%#v ok=%v", restored, ok)
	}
	if !governance.Verify(second.audit.List()) || len(second.audit.List()) != 1 {
		t.Fatalf("audit=%#v", second.audit.List())
	}
}
