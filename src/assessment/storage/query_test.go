package storage

import (
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
)

func TestQueryJobsFiltersAndPaginates(t *testing.T) {
	now := time.Now().UTC()
	jobs := []domain.AssessmentJob{
		{ID: "1", EngagementID: "eng", Status: "COMPLETED", Target: "fixture://a", CreatedAt: now},
		{ID: "2", EngagementID: "eng", Status: "FAILED", Target: "fixture://b", CreatedAt: now.Add(-time.Minute)},
		{ID: "3", EngagementID: "other", Status: "COMPLETED", Target: "fixture://a", CreatedAt: now.Add(-2 * time.Minute)},
	}
	page := QueryJobs(jobs, JobQuery{EngagementID: "eng", Status: "completed", Page: 1, PageSize: 1})
	if page.Total != 1 || len(page.Items) != 1 || page.Items[0].ID != "1" || page.HasNext {
		t.Fatalf("page=%#v", page)
	}
}
