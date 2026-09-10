package storage

import (
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
)

func TestQueryObservationsFiltersByJobAndStatus(t *testing.T) {
	now := time.Now().UTC()
	items := []domain.Observation{{ID: "1", JobID: "job-1", CheckID: "check", Status: "PASSED", CreatedAt: now}, {ID: "2", JobID: "job-2", CheckID: "check", Status: "FAILED", CreatedAt: now.Add(-time.Minute)}}
	page := QueryObservations(items, ObservationQuery{JobID: "job-1", Status: "passed", Page: 1, PageSize: 25})
	if page.Total != 1 || len(page.Items) != 1 || page.Items[0].ID != "1" {
		t.Fatalf("page=%#v", page)
	}
}
