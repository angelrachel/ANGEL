package evidence

import "testing"

func TestManifestRoundTrip(t *testing.T) {
	ledger := NewLedger()
	ledger.AddRecord(`{"event":"start"}`)
	ledger.AddRecord(`{"event":"finish"}`)
	manifest, err := BuildManifest(ledger.Records)
	if err != nil {
		t.Fatalf("build manifest: %v", err)
	}
	if err := manifest.Validate(ledger.Records); err != nil {
		t.Fatalf("validate manifest: %v", err)
	}
	if manifest.RecordCount != 2 || manifest.FirstID == "" || manifest.LastID == "" || manifest.ChainHash == "" {
		t.Fatalf("incomplete manifest: %#v", manifest)
	}
}

func TestManifestRejectsTamperedChain(t *testing.T) {
	ledger := NewLedger()
	ledger.AddRecord("event")
	records := append([]Record(nil), ledger.Records...)
	manifest, err := BuildManifest(records)
	if err != nil {
		t.Fatalf("build manifest: %v", err)
	}
	records[0].Data = "tampered"
	if err := manifest.Validate(records); err == nil {
		t.Fatal("expected tampered records to be rejected")
	}
}
