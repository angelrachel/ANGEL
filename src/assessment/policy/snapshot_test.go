package policy

import (
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
)

func TestSnapshotRestorePreservesBudgetUsageAndScope(t *testing.T) {
	now := time.Now().UTC()
	first := NewEngine()
	engagement := domain.Engagement{ID: "eng-snapshot", Organization: "org", Authorized: true, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Hour)}
	scope := []domain.ScopeEntry{{Kind: "fixture", Value: "fixture://lab/web", Actions: []string{"surface-map"}}}
	if err := first.RegisterEngagement(engagement, scope, 2); err != nil {
		t.Fatal(err)
	}
	if !first.Authorize(engagement.ID, "fixture://lab/web", "surface-map", now).Allowed {
		t.Fatal("initial authorization denied")
	}
	second := NewEngine()
	if err := second.Restore(first.Snapshot()); err != nil {
		t.Fatal(err)
	}
	if !second.Authorize(engagement.ID, "fixture://lab/web", "surface-map", now).Allowed {
		t.Fatal("restored authorization denied")
	}
	if second.Authorize(engagement.ID, "fixture://lab/web", "surface-map", now).Allowed {
		t.Fatal("restored budget usage lost")
	}
}
