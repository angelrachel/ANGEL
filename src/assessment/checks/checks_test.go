package checks

import (
	"crypto/tls"
	"testing"
)

func TestRegistryContainsConcreteChecks(t *testing.T) {
	if len(Registry()) < 12 {
		t.Fatalf("checks=%d", len(Registry()))
	}
	if _, ok := Find("web-security-headers"); !ok {
		t.Fatal("web check missing")
	}
	if _, ok := Find("evidence-integrity-check"); !ok {
		t.Fatal("evidence check missing")
	}
}

func TestSecurityHeadersAndSecretExposure(t *testing.T) {
	missing := SecurityHeaders{}.Evaluate(Request{Headers: map[string]string{}})
	if missing.Passed || missing.Severity != High {
		t.Fatalf("unexpected header result: %#v", missing)
	}
	exposed := SecretExposure{}.Evaluate(Request{Body: "authorization: bearer abc"})
	if exposed.Passed || exposed.Severity != Critical {
		t.Fatalf("unexpected secret result: %#v", exposed)
	}
}

func TestTLSPostureRequiresModernProtocol(t *testing.T) {
	result := TLSPosture{}.Evaluate(Request{TLSState: &tls.ConnectionState{Version: tls.VersionTLS12, CipherSuite: tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256}})
	if result.Passed || result.Severity != High {
		t.Fatalf("unexpected TLS result: %#v", result)
	}
}

func TestAPIContractAndEvidenceIntegrity(t *testing.T) {
	valid := APISchema{}.Evaluate(Request{Body: `{"status":"ok"}`})
	if !valid.Passed {
		t.Fatalf("valid JSON rejected: %#v", valid)
	}
	invalid := EvidenceIntegrity{}.Evaluate(Request{Metadata: map[string]string{"sha256": "bad"}})
	if invalid.Passed {
		t.Fatal("invalid evidence hash accepted")
	}
}
