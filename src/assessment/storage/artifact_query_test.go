package storage

import (
	"testing"
	"time"

	"ANGEL/src/assessment/evidence"
	"ANGEL/src/assessment/reporting"
)

func TestQueryEvidenceRedactsPayloadAndFiltersJob(t *testing.T) {
	page := QueryEvidence([]evidence.Bundle{{ID: "e1", JobID: "job-1", EngagementID: "eng", ContentType: "application/json", Payload: []byte("secret"), CapturedAt: time.Now()}, {ID: "e2", JobID: "job-2", EngagementID: "eng", CapturedAt: time.Now().Add(-time.Minute)}}, EvidenceQuery{JobID: "job-1"})
	if page.Total != 1 || len(page.Items) != 1 || len(page.Items[0].Payload) != 0 {
		t.Fatalf("page=%#v", page)
	}
}

func TestQueryReportsFiltersEngagement(t *testing.T) {
	page := QueryReports([]reporting.Report{{ID: "r1", EngagementID: "eng", GeneratedAt: time.Now()}, {ID: "r2", EngagementID: "other", GeneratedAt: time.Now()}}, ReportQuery{EngagementID: "eng"})
	if page.Total != 1 || len(page.Items) != 1 || page.Items[0].ID != "r1" {
		t.Fatalf("page=%#v", page)
	}
}
