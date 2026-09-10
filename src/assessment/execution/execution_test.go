package execution

import (
	"context"
	"strings"
	"testing"
)

func TestFixtureAdapterProducesEvidenceAndFinding(t *testing.T) {
	output, err := NewRegistry().Run(context.Background(), Input{JobID: "job-1", Target: "fixture://lab/web", CheckID: "web-security-headers", Fixture: map[string]any{"body": "ok"}})
	if err != nil {
		t.Fatal(err)
	}
	if output.Status != "COMPLETED" || len(output.Evidence) == 0 || output.Result.CheckID != "web-security-headers" {
		t.Fatalf("output=%#v", output)
	}
}
func TestRegistryRejectsUnsupportedCheck(t *testing.T) {
	_, err := NewRegistry().Run(context.Background(), Input{Target: "fixture://lab", CheckID: "unknown"})
	if err == nil || !strings.Contains(err.Error(), "no safe adapter") {
		t.Fatalf("err=%v", err)
	}
}
func TestHTTPAdapterRejectsFixtureTarget(t *testing.T) {
	_, err := (HTTPMetadataAdapter{}).Run(context.Background(), Input{Target: "fixture://lab", CheckID: "web-method-policy"})
	if err == nil {
		t.Fatal("fixture accepted by HTTP adapter")
	}
}
