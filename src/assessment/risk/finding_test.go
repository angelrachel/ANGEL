package risk

import "testing"

func TestFindingRequiresEvidenceImpactAndRemediation(t *testing.T) {
	input := Input{Exploitability: 8, Reachability: 8, Confidentiality: 8, Integrity: 7, Availability: 6, BusinessCriticality: 8, BlastRadius: 7, Confidence: .92, DetectionCoverage: 2}
	if _, _, err := Finding("title", "impact", "remediation", []string{"asset-1"}, []string{"evidence-1"}, input); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Finding("title", "", "remediation", []string{"asset-1"}, []string{"evidence-1"}, input); err == nil {
		t.Fatal("missing impact accepted")
	}
	if _, _, err := Finding("title", "impact", "remediation", nil, []string{"evidence-1"}, input); err == nil {
		t.Fatal("missing asset accepted")
	}
}

func TestFindingIdentityIsDeterministic(t *testing.T) {
	input := Input{Exploitability: 5, Reachability: 5, Confidentiality: 5, Integrity: 5, Availability: 5, BusinessCriticality: 5, BlastRadius: 5, Confidence: .8, DetectionCoverage: 5}
	first, _, err := Finding("same", "impact", "fix", []string{"asset"}, []string{"evidence"}, input)
	if err != nil {
		t.Fatal(err)
	}
	second, _, err := Finding("same", "impact", "fix", []string{"asset"}, []string{"evidence"}, input)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID == "" || first.ID != second.ID {
		t.Fatalf("identity is not deterministic: %q %q", first.ID, second.ID)
	}
}
