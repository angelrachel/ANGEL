package report

import (
	"strings"
	"testing"
)

func TestJSONExportSortsFindingsBySeverity(t *testing.T) {
	report := GenerateReport("Assessment", []Finding{{Title: "Low", Severity: "low", Description: "d", Evidence: "e"}, {Title: "Critical", Severity: "critical", Description: "d", Evidence: "e"}})
	value, err := report.JSON()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Index(value, "Critical") > strings.Index(value, "Low") {
		t.Fatalf("findings not severity ordered: %s", value)
	}
}

func TestJSONExportRejectsInvalidReport(t *testing.T) {
	if _, err := (Report{ID: "x", Title: "bad", StartTime: "invalid", EndTime: "invalid"}).JSON(); err == nil {
		t.Fatal("expected invalid report rejection")
	}
}

func TestExecutiveSummaryCountsSeverity(t *testing.T) {
	report := Report{Findings: []Finding{{Severity: "critical"}, {Severity: "high"}, {Severity: "high"}}}
	if got := report.ExecutiveSummary(); got != "3 finding(s): 1 critical, 2 high, 0 medium, 0 low, 0 info." {
		t.Fatalf("unexpected summary: %s", got)
	}
}
