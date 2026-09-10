package evidence

import (
	"testing"
	"time"
)

func TestBundleIsSignedAndTamperEvident(t *testing.T) {
	signer, err := NewSigner()
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := CreateBundle("eng", "job", "application/json", []byte(`{"status":"ok"}`), "0", true, time.Now(), signer)
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyBundle(bundle, signer.Public) {
		t.Fatal("bundle did not verify")
	}
	bundle.Payload[0] = 'x'
	if VerifyBundle(bundle, signer.Public) {
		t.Fatal("tampered bundle verified")
	}
}
