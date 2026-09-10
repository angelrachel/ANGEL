package persistence

import (
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
)

func TestSnapshotRoundTripAndTamperDetection(t *testing.T) {
	signer, err := NewSigner()
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := Create(Snapshot{Version: 1, Engagements: []domain.Engagement{{ID: "eng-1", Authorized: true}}, ExportedAt: time.Now().UTC()}, signer)
	if err != nil {
		t.Fatal(err)
	}
	if err := Verify(snapshot, signer.Public); err != nil {
		t.Fatal(err)
	}
	raw, err := Encode(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if err := Verify(decoded, signer.Public); err != nil {
		t.Fatal(err)
	}
	decoded.Engagements[0].ID = "tampered"
	if Verify(decoded, signer.Public) == nil {
		t.Fatal("tampered snapshot verified")
	}
}
func TestUnsignedSnapshotCannotBeEncoded(t *testing.T) {
	if _, err := Encode(Snapshot{Version: 1}); err == nil {
		t.Fatal("unsigned snapshot encoded")
	}
}
