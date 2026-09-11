package reporting

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateReportCreatesDeterministicID(t *testing.T) {
	findings := []Finding{
		{Title: "Test Finding", Severity: "P1", Description: "Test description", Evidence: "ev-1"},
	}
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	report1 := GenerateReportAt("Test Report", findings, fixedTime)
	report2 := GenerateReportAt("Test Report", findings, fixedTime)

	if report1.ID != report2.ID {
		t.Fatalf("expected deterministic ID, got %q vs %q", report1.ID, report2.ID)
	}
	if len(report1.ID) != 64 {
		t.Fatalf("expected SHA256 hex (64 chars), got %d", len(report1.ID))
	}
}

func TestGenerateReportIncludesTimestamp(t *testing.T) {
	findings := []Finding{}
	report := GenerateReport("Test", findings)
	if report.StartTime == "" {
		t.Fatal("expected StartTime")
	}
	if report.EndTime == "" {
		t.Fatal("expected EndTime")
	}
	if report.StartTime != report.EndTime {
		t.Fatalf("expected same start/end for instant report: %q vs %q", report.StartTime, report.EndTime)
	}
}

func TestGenerateReportWithMultipleFindings(t *testing.T) {
	findings := []Finding{
		{Title: "F1", Severity: "P2", Description: "D1", Evidence: "E1"},
		{Title: "F2", Severity: "P1", Description: "D2", Evidence: "E2"},
	}
	report := GenerateReport("Multi", findings)
	if len(report.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(report.Findings))
	}
	if report.Findings[0].Title != "F1" {
		t.Fatalf("expected F1, got %q", report.Findings[0].Title)
	}
	if report.Findings[1].Title != "F2" {
		t.Fatalf("expected F2, got %q", report.Findings[1].Title)
	}
}

func TestReportMarkdownOutput(t *testing.T) {
	findings := []Finding{
		{Title: "XSS", Severity: "P1", Description: "Reflected XSS found", Evidence: "ev-123"},
	}
	report := GenerateReport("Test", findings)
	md := report.Markdown()

	if !strings.Contains(md, "# Test") {
		t.Fatal("expected title in markdown")
	}
	if !strings.Contains(md, "## Findings") {
		t.Fatal("expected findings section")
	}
	if !strings.Contains(md, "### XSS") {
		t.Fatal("expected finding title")
	}
	if !strings.Contains(md, "**Severity:** P1") {
		t.Fatal("expected severity")
	}
	if !strings.Contains(md, "**Description:** Reflected XSS found") {
		t.Fatal("expected description")
	}
	if !strings.Contains(md, "**Evidence:** ev-123") {
		t.Fatal("expected evidence")
	}
}

func TestReportMarkdownEmptyFindings(t *testing.T) {
	report := GenerateReport("Empty", nil)
	md := report.Markdown()
	if !strings.Contains(md, "# Empty") {
		t.Fatal("expected title")
	}
	if !strings.Contains(md, "## Findings") {
		t.Fatal("expected findings section even when empty")
	}
}
