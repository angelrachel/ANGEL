package evidence

import (
	"bytes"
	"testing"
)

func TestLedgerExportImportPreservesHead(t *testing.T) {
	ledger := NewLedger()
	ledger.AddRecord(`{"type":"observation","id":"1"}`)
	ledger.AddRecord(`{"type":"finding","id":"2"}`)
	exported, err := ledger.Export()
	if err != nil {
		t.Fatal(err)
	}
	restored, err := Import(exported)
	if err != nil {
		t.Fatal(err)
	}
	head, ok := restored.Head()
	if !ok || head.Hash != ledger.Records[1].Hash || !VerifyChain(restored.Records) {
		t.Fatalf("restored=%#v", restored)
	}
}

func TestLedgerImportRejectsTampering(t *testing.T) {
	ledger := NewLedger()
	ledger.AddRecord("safe")
	exported, _ := ledger.Export()
	tampered := bytes.Replace(exported, []byte("safe"), []byte("changed"), 1)
	if _, err := Import(tampered); err == nil {
		t.Fatal("tampered ledger accepted")
	}
}
