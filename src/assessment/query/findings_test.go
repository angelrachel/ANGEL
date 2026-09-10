package query

import (
	"testing"

	"ANGEL/src/assessment/reporting"
)

func TestFindingsFilterSortsAndPaginates(t *testing.T) {
	items := []reporting.Finding{{ID: "low", CheckID: "tls", Title: "TLS", Severity: "LOW", Confidence: .9, Target: "a", Status: "OPEN"}, {ID: "p0", CheckID: "secret", Title: "Secret exposure", Severity: "CRITICAL", Confidence: .99, Target: "a", Status: "OPEN"}, {ID: "p1", CheckID: "header", Title: "Headers", Severity: "HIGH", Confidence: .8, Target: "b", Status: "ACCEPTED"}}
	page := Findings(items, FindingFilter{Target: "a", Search: "secret", Page: 1, PageSize: 1})
	if page.Total != 1 || len(page.Items) != 1 || page.Items[0].ID != "p0" || page.HasNext {
		t.Fatalf("page=%#v", page)
	}
}
func TestFindingsPaginationCapsPageSize(t *testing.T) {
	page := Findings(nil, FindingFilter{Page: 0, PageSize: 1000})
	if page.Page != 1 || page.PageSize != 100 || page.Items == nil {
		t.Fatalf("page=%#v", page)
	}
}
