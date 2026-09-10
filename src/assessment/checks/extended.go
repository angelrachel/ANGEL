package checks

import (
	"fmt"
	"net/url"
	"strings"
)

type AuthenticationPosture struct{}

func (AuthenticationPosture) ID() string { return "authentication-session-posture" }
func (AuthenticationPosture) Layer() int { return 11 }
func (AuthenticationPosture) Evaluate(r Request) Result {
	rotation := strings.EqualFold(r.Metadata["refresh_rotation"], "true")
	mfa := strings.EqualFold(r.Metadata["mfa_required"], "true")
	passed := rotation && mfa
	return result("authentication-session-posture", 11, "Authentication session posture", passed, High, .95, fmt.Sprintf("refresh rotation=%t, mfa=%t", rotation, mfa), map[string]string{"refresh_rotation": fmt.Sprint(rotation), "mfa_required": fmt.Sprint(mfa)}, "Require MFA for privileged roles and rotate refresh tokens on use")
}

type TenantIsolation struct{}

func (TenantIsolation) ID() string { return "authorization-tenant-isolation" }
func (TenantIsolation) Layer() int { return 12 }
func (TenantIsolation) Evaluate(r Request) Result {
	crossTenant := strings.EqualFold(r.Metadata["cross_tenant_access"], "true")
	return result("authorization-tenant-isolation", 12, "Tenant isolation", !crossTenant, Critical, .98, "cross-tenant access simulation evaluated", map[string]string{"cross_tenant_access": fmt.Sprint(crossTenant)}, "Enforce tenant identity from the authenticated session and deny cross-tenant object access")
}

type SSRFCanary struct{}

func (SSRFCanary) ID() string { return "server-side-canary-policy" }
func (SSRFCanary) Layer() int { return 14 }
func (SSRFCanary) Evaluate(r Request) Result {
	parsed, err := url.Parse(r.Target)
	valid := err == nil && parsed.Hostname() != "" && (parsed.Scheme == "https" || parsed.Scheme == "http")
	callback := strings.EqualFold(r.Metadata["canary_observed"], "true")
	passed := valid && !callback
	return result("server-side-canary-policy", 14, "Server-side request canary policy", passed, High, .94, fmt.Sprintf("url_valid=%t, canary_observed=%t", valid, callback), map[string]string{"url_valid": fmt.Sprint(valid), "canary_observed": fmt.Sprint(callback)}, "Allow only approved destinations and use synthetic callback markers without fetching internal data")
}

type APIRateLimit struct{}

func (APIRateLimit) ID() string { return "api-rate-limit-enforcement" }
func (APIRateLimit) Layer() int { return 15 }
func (APIRateLimit) Evaluate(r Request) Result {
	enforced := strings.EqualFold(r.Metadata["rate_limit_enforced"], "true")
	return result("api-rate-limit-enforcement", 15, "API rate limit enforcement", enforced, High, .93, "bounded request budget evaluated", map[string]string{"rate_limit_enforced": fmt.Sprint(enforced)}, "Enforce rate and concurrency limits per identity, engagement, target, and endpoint")
}

type PrivilegePath struct{}

func (PrivilegePath) ID() string { return "identity-privilege-path" }
func (PrivilegePath) Layer() int { return 16 }
func (PrivilegePath) Evaluate(r Request) Result {
	reachable := strings.EqualFold(r.Metadata["production_reachable"], "true")
	excessive := strings.EqualFold(r.Metadata["excessive_permission"], "true")
	passed := !reachable && !excessive
	return result("identity-privilege-path", 16, "Identity privilege path", passed, Critical, .95, fmt.Sprintf("production_reachable=%t, excessive_permission=%t", reachable, excessive), map[string]string{"production_reachable": fmt.Sprint(reachable), "excessive_permission": fmt.Sprint(excessive)}, "Remove unnecessary privilege and require explicit, reviewable transitions to critical assets")
}

type DetectionLatency struct{}

func (DetectionLatency) ID() string { return "detection-alert-latency" }
func (DetectionLatency) Layer() int { return 19 }
func (DetectionLatency) Evaluate(r Request) Result {
	detected := strings.EqualFold(r.Metadata["detected"], "true")
	latency := r.Metadata["latency_seconds"]
	passed := detected && latency != ""
	return result("detection-alert-latency", 19, "Detection alert latency", passed, High, .9, fmt.Sprintf("detected=%t, latency_seconds=%s", detected, latency), map[string]string{"detected": fmt.Sprint(detected), "latency_seconds": latency}, "Generate a benign marker, verify ingestion and alerting, and define an approved response-time objective")
}

type SupplyChainProvenance struct{}

func (SupplyChainProvenance) ID() string { return "supply-chain-provenance" }
func (SupplyChainProvenance) Layer() int { return 20 }
func (SupplyChainProvenance) Evaluate(r Request) Result {
	sbom := strings.EqualFold(r.Metadata["sbom_present"], "true")
	signed := strings.EqualFold(r.Metadata["provenance_verified"], "true")
	passed := sbom && signed
	return result("supply-chain-provenance", 20, "Supply-chain provenance", passed, High, .94, fmt.Sprintf("sbom_present=%t, provenance_verified=%t", sbom, signed), map[string]string{"sbom_present": fmt.Sprint(sbom), "provenance_verified": fmt.Sprint(signed)}, "Publish an accurate SBOM and verify signed build provenance before release")
}

type RecoveryMeasurement struct{}

func (RecoveryMeasurement) ID() string { return "recovery-rto-rpo-measurement" }
func (RecoveryMeasurement) Layer() int { return 21 }
func (RecoveryMeasurement) Evaluate(r Request) Result {
	rto := strings.TrimSpace(r.Metadata["rto_seconds"])
	rpo := strings.TrimSpace(r.Metadata["rpo_seconds"])
	restored := strings.EqualFold(r.Metadata["restore_verified"], "true")
	passed := restored && rto != "" && rpo != ""
	return result("recovery-rto-rpo-measurement", 21, "Recovery RTO/RPO measurement", passed, High, .92, fmt.Sprintf("restore_verified=%t, rto=%s, rpo=%s", restored, rto, rpo), map[string]string{"restore_verified": fmt.Sprint(restored), "rto_seconds": rto, "rpo_seconds": rpo}, "Run a disposable restore test and record measured RTO/RPO against approved objectives")
}
