package control

import (
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
	"ANGEL/src/assessment/policy"
)

func TestControllerCreatesSignedJobAndTransitions(t *testing.T) {
	now := time.Now().UTC()
	p := policy.NewEngine()
	if err := p.RegisterEngagement(domain.Engagement{ID: "eng-1", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"surface-map"}}}, 4); err != nil {
		t.Fatal(err)
	}
	controller, err := NewController(p)
	if err != nil {
		t.Fatal(err)
	}
	job, err := controller.Create("eng-1", "surface-map", "fixture://lab/web", "surface-map", now)
	if err != nil {
		t.Fatal(err)
	}
	if job.Signature == "" {
		t.Fatal("job signature missing")
	}
	for _, state := range []string{Approved, Queued, Running, Completed} {
		if err := controller.Transition(job.ID, state, now); err != nil {
			t.Fatalf("transition %s: %v", state, err)
		}
	}
	final, _ := controller.Get(job.ID)
	if final.Status != Completed {
		t.Fatalf("status=%s", final.Status)
	}
}

func TestControllerRejectsIllegalTransition(t *testing.T) {
	now := time.Now().UTC()
	p := policy.NewEngine()
	_ = p.RegisterEngagement(domain.Engagement{ID: "eng-2", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}, []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"surface-map"}}}, 1)
	controller, _ := NewController(p)
	job, err := controller.Create("eng-2", "surface-map", "fixture://lab/web", "surface-map", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := controller.Transition(job.ID, Completed, now); err == nil {
		t.Fatal("illegal transition accepted")
	}
}
