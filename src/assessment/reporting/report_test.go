package reporting

import (
	"strings"
	"testing"
	"time"

	"ANGEL/src/assessment/checks"
)

func TestReportIsValidatedAndRendered(t *testing.T) {
	finding := FromCheck(checks.SecurityHeaders{}.Evaluate(checks.Request{Headers: map[string]string{}}), "fixture://lab/web", "evidence-1")
	report, err := New("eng-1", []Finding{finding}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := report.Validate(); err != nil {
		t.Fatal(err)
	}
	raw, err := report.JSON()
	if err != nil || len(raw) == 0 {
		t.Fatalf("json err=%v", err)
	}
	if !strings.Contains(report.Markdown(), "Security Validation Report") {
		t.Fatal("markdown title missing")
	}
}
