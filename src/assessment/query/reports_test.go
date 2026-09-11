package query

import (
	"testing"
	"time"

	"ANGEL/src/assessment/reporting"
)

func TestReportsFiltersSortsAndPaginates(t *testing.T) {
	now := time.Now().UTC()
	items := []reporting.Report{
		{ID: "old", EngagementID: "eng-1", GeneratedAt: now.Add(-time.Hour)},
		{ID: "new", EngagementID: "eng-1", GeneratedAt: now},
		{ID: "other", EngagementID: "eng-2", GeneratedAt: now.Add(time.Hour)},
	}
	page := Reports(items, ReportFilter{EngagementID: "eng-1", Page: 1, PageSize: 1})
	if page.Total != 2 || len(page.Items) != 1 || page.Items[0].ID != "new" || !page.HasNext {
		t.Fatalf("unexpected page: %+v", page)
	}
}
