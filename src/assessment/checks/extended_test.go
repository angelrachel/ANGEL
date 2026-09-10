package checks

import "testing"

func TestExtendedRegistryChecks(t *testing.T) {
	for _, id := range []string{
		"authentication-session-posture", "authorization-tenant-isolation", "server-side-canary-policy",
		"api-rate-limit-enforcement", "identity-privilege-path", "detection-alert-latency",
		"supply-chain-provenance", "recovery-rto-rpo-measurement",
	} {
		if _, ok := Find(id); !ok {
			t.Fatalf("missing extended check %s", id)
		}
	}
}

func TestExtendedChecksRejectUnsafePosture(t *testing.T) {
	if got := (TenantIsolation{}).Evaluate(Request{Metadata: map[string]string{"cross_tenant_access": "true"}}); got.Passed || got.Severity != Critical {
		t.Fatalf("tenant isolation unsafe posture accepted: %#v", got)
	}
	if got := (SSRFCanary{}).Evaluate(Request{Target: "http://fixture.invalid", Metadata: map[string]string{"canary_observed": "true"}}); got.Passed {
		t.Fatalf("observed canary accepted: %#v", got)
	}
	if got := (SupplyChainProvenance{}).Evaluate(Request{Metadata: map[string]string{"sbom_present": "true"}}); got.Passed {
		t.Fatalf("unsigned provenance accepted: %#v", got)
	}
}

func TestExtendedChecksAcceptVerifiedPosture(t *testing.T) {
	metadata := map[string]string{"refresh_rotation": "true", "mfa_required": "true"}
	if got := (AuthenticationPosture{}).Evaluate(Request{Metadata: metadata}); !got.Passed {
		t.Fatalf("verified authentication posture rejected: %#v", got)
	}
	if got := (RecoveryMeasurement{}).Evaluate(Request{Metadata: map[string]string{"restore_verified": "true", "rto_seconds": "30", "rpo_seconds": "10"}}); !got.Passed {
		t.Fatalf("measured recovery posture rejected: %#v", got)
	}
}
