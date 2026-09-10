package risk

import "testing"

func TestScoreProducesP0ForHighConfidenceCriticalPath(t *testing.T) {
	result, err := Score(Input{Exploitability: 10, Reachability: 10, Confidentiality: 10, Integrity: 10, Availability: 9, BusinessCriticality: 10, BlastRadius: 10, Confidence: .98, DetectionCoverage: 0})
	if err != nil {
		t.Fatal(err)
	}
	if result.Severity != "P0" {
		t.Fatalf("severity=%s score=%.2f", result.Severity, result.Score)
	}
}

func TestScoreRejectsInvalidConfidence(t *testing.T) {
	if _, err := Score(Input{Confidence: 2}); err == nil {
		t.Fatal("invalid confidence accepted")
	}
}
