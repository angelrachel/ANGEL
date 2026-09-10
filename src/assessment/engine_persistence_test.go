package main

import (
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
)

func TestEngineRestoresDurableSignedJobState(t *testing.T) {
	path := t.TempDir() + "/angel-state.json"
	t.Setenv("ANGEL_STATE_PATH", path)
	now := time.Now().UTC()
	first := NewEngine()
	engagement := domain.Engagement{ID: "persist-eng", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Hour)}
	if err := first.policyEngine.RegisterEngagement(engagement, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"surface-map"}}}, 2); err != nil {
		t.Fatal(err)
	}
	job, err := first.jobController.Create(engagement.ID, "surface-map", "fixture://lab/web", "surface-map", now)
	if err != nil {
		t.Fatal(err)
	}
	first.persistState()
	second := NewEngine()
	restored, ok := second.jobController.Get(job.ID)
	if !ok || restored.ID != job.ID || !second.jobController.Verify(restored) {
		t.Fatalf("restored=%#v ok=%v", restored, ok)
	}
}
