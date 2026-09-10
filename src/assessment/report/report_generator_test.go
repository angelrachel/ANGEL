package report

import (
	"strings"
	"testing"
)

func TestGeneratedReportValidates(t *testing.T) {
	report := GenerateReport("Assessment", []Finding{{
		Title: "Missing header", Severity: "medium", Description: "Header is absent", Evidence: "evidence-1",
	}})
	if err := report.Validate(); err != nil {
		t.Fatalf("expected generated report to validate: %v", err)
	}
}

func TestValidateRejectsIncompleteFindingAndBadSeverity(t *testing.T) {
	report := GenerateReport("Assessment", []Finding{{
		Title: "", Severity: "medium", Description: "description",
	}})
	if err := report.Validate(); err == nil {
		t.Fatal("expected empty finding title to be rejected")
	}

	report = GenerateReport("Assessment", []Finding{{
		Title: "Finding", Severity: "urgent", Description: "description",
	}})
	if err := report.Validate(); err == nil {
		t.Fatal("expected unsupported severity to be rejected")
	}
}

func TestValidateRejectsReversedTimes(t *testing.T) {
	report := GenerateReport("Assessment", nil)
	report.StartTime = "2026-09-10T01:00:00Z"
	report.EndTime = "2026-09-10T00:00:00Z"
	if err := report.Validate(); err == nil {
		t.Fatal("expected reversed report times to be rejected")
	}
}

func TestMarkdownEscapesControlMarkup(t *testing.T) {
	report := Report{
		Title: "Title #1",
		Findings: []Finding{{
			Title: "[untrusted] *title*", Severity: "high", Description: "Use `code` and _emphasis_", Evidence: "evidence",
		}},
	}
	markdown := report.Markdown()
	for _, token := range []string{"\\#", "\\[", "\\*", "\\`", "\\_"} {
		if !strings.Contains(markdown, token) {
			t.Fatalf("expected escaped token %q in %q", token, markdown)
		}
	}
}

func TestNilReportMarkdownIsEmpty(t *testing.T) {
	var report *Report
	if report.Markdown() != "" {
		t.Fatal("expected nil report markdown to be empty")
	}
}
