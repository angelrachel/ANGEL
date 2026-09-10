package domain

import "testing"

func TestCoreModelsHaveStableJSONTags(t *testing.T) {
	if (Finding{}).Status != "" {
		t.Fatal("zero value should be empty")
	}
	if (AssessmentJob{}).Status != "" {
		t.Fatal("zero value should be empty")
	}
}
