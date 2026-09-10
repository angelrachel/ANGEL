package report

import "testing"

func TestRenderArtifactFormatsAndDigest(t *testing.T) {
	report := GenerateReport("ANGEL assessment", []Finding{{Title: "Header gap", Severity: "HIGH", Description: "CSP missing", Evidence: "bundle-1"}})
	for _, format := range []string{"json", "markdown", "html"} {
		artifact, err := RenderArtifact(report, format)
		if err != nil {
			t.Fatalf("format %s: %v", format, err)
		}
		if len(artifact.Body) == 0 || len(artifact.SHA256) != 64 || artifact.ContentType == "" {
			t.Fatalf("invalid artifact %#v", artifact)
		}
	}
}

func TestRenderArtifactRejectsInvalidReportAndFormat(t *testing.T) {
	if _, err := RenderArtifact(Report{}, "json"); err == nil {
		t.Fatal("invalid report accepted")
	}
	report := GenerateReport("ANGEL", nil)
	if _, err := RenderArtifact(report, "pdf"); err == nil {
		t.Fatal("unsupported format accepted")
	}
}
