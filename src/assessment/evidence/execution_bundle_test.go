package evidence

import (
	"testing"
	"time"

	"ANGEL/src/assessment/checks"
	"ANGEL/src/assessment/collectors"
	"ANGEL/src/assessment/execution"
)

func TestBundleExecutionCreatesVerifiableEvidence(t *testing.T) {
	signer, err := NewSigner()
	if err != nil {
		t.Fatal(err)
	}
	output := execution.Output{JobID: "job-1", CheckID: "web-security-headers", Status: "COMPLETED", Result: checks.Result{CheckID: "web-security-headers", Passed: true}, Observation: collectors.Observation{Kind: "http", Target: "fixture://http/security"}, Evidence: []byte(`{"fixture":true}`), CollectedAt: time.Now().UTC()}
	bundle, err := BundleExecution(output, "eng-1", "parent-hash", signer, time.Now().UTC())
	if err != nil {
		t.Fatalf("bundle execution: %v", err)
	}
	if bundle.ParentSHA256 != "parent-hash" || !bundle.Redacted || bundle.ContentType == "" {
		t.Fatalf("unexpected bundle: %+v", bundle)
	}
	if !VerifyBundle(bundle, signer.Public) {
		t.Fatal("bundle signature verification failed")
	}
}

func TestBundleExecutionRejectsIncompleteOutput(t *testing.T) {
	signer, err := NewSigner()
	if err != nil {
		t.Fatal(err)
	}
	_, err = BundleExecution(execution.Output{JobID: "job-1", CheckID: "check-1", Status: "RUNNING"}, "eng-1", "", signer, time.Now())
	if err == nil {
		t.Fatal("expected incomplete execution to be rejected")
	}
}
