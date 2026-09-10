package normalization

import (
	"testing"
	"time"

	"ANGEL/src/assessment/collectors"
)

func TestCanonicalAsset(t *testing.T) {
	asset, err := New().Asset("HTTPS://Example.TEST/path", "web", []string{"External", "external"}, .8)
	if err != nil {
		t.Fatal(err)
	}
	if asset.Canonical != "https://example.test" || len(asset.Tags) != 1 {
		t.Fatalf("asset=%#v", asset)
	}
}
func TestMergeObservationRaisesConfidence(t *testing.T) {
	asset := Asset{Canonical: "fixture://lab/web", Confidence: .5, LastSeen: time.Now()}
	merged := MergeObservation(asset, collectors.Observation{Target: "fixture://lab/web", Kind: "http", Confidence: .9, CollectedAt: time.Now()})
	if merged.Confidence != .9 || len(merged.Tags) != 1 {
		t.Fatalf("merged=%#v", merged)
	}
}
