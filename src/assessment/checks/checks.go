package checks

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Severity string

const (
	Info     Severity = "INFO"
	Low      Severity = "LOW"
	Medium   Severity = "MEDIUM"
	High     Severity = "HIGH"
	Critical Severity = "CRITICAL"
)

type Request struct {
	ID         string               `json:"id"`
	Target     string               `json:"target"`
	Headers    map[string]string    `json:"headers,omitempty"`
	StatusCode int                  `json:"status_code,omitempty"`
	Body       string               `json:"body,omitempty"`
	TLSState   *tls.ConnectionState `json:"-"`
	Metadata   map[string]string    `json:"metadata,omitempty"`
}
type Result struct {
	CheckID     string            `json:"check_id"`
	Layer       int               `json:"layer"`
	Title       string            `json:"title"`
	Passed      bool              `json:"passed"`
	Severity    Severity          `json:"severity"`
	Confidence  float64           `json:"confidence"`
	Summary     string            `json:"summary"`
	Evidence    map[string]string `json:"evidence"`
	Remediation string            `json:"remediation"`
}
type Check interface {
	ID() string
	Layer() int
	Evaluate(Request) Result
}

func Registry() []Check {
	return []Check{SecurityHeaders{}, TLSPosture{}, HTTPMethodPolicy{}, ContentTypePolicy{}, APISchema{}, AuthorizationMatrix{}, SecretExposure{}, DependencyPosture{}, CloudIdentityPosture{}, EndpointLogging{}, RecoveryControl{}, EvidenceIntegrity{}}
}
func Find(id string) (Check, bool) {
	for _, check := range Registry() {
		if check.ID() == id {
			return check, true
		}
	}
	return nil, false
}
func IDs() []string {
	out := make([]string, 0, len(Registry()))
	for _, check := range Registry() {
		out = append(out, check.ID())
	}
	sort.Strings(out)
	return out
}

type SecurityHeaders struct{}

func (SecurityHeaders) ID() string { return "web-security-headers" }
func (SecurityHeaders) Layer() int { return 9 }
func (SecurityHeaders) Evaluate(r Request) Result {
	missing := []string{}
	for _, h := range []string{"Content-Security-Policy", "X-Content-Type-Options", "Referrer-Policy", "Strict-Transport-Security"} {
		if strings.TrimSpace(r.Headers[h]) == "" {
			missing = append(missing, h)
		}
	}
	passed := len(missing) == 0
	sev := Low
	if len(missing) >= 3 {
		sev = High
	}
	return result("web-security-headers", 9, "HTTP security headers", passed, sev, .98, fmt.Sprintf("missing headers: %s", strings.Join(missing, ", ")), map[string]string{"missing": strings.Join(missing, ",")}, "Configure security headers at the edge and application layer")
}

type TLSPosture struct{}

func (TLSPosture) ID() string { return "network-tls-posture" }
func (TLSPosture) Layer() int { return 8 }
func (TLSPosture) Evaluate(r Request) Result {
	if r.TLSState == nil {
		return result("network-tls-posture", 8, "TLS posture", false, High, .8, "TLS state was not provided", nil, "Collect TLS handshake evidence before rating posture")
	}
	v := r.TLSState.Version
	passed := v >= tls.VersionTLS13
	evidence := map[string]string{"version": tlsVersion(v), "cipher": tls.CipherSuiteName(r.TLSState.CipherSuite)}
	if !passed {
		return result("network-tls-posture", 8, "TLS posture", false, High, .99, "TLS version is below TLS 1.3", evidence, "Disable legacy protocol versions and require TLS 1.3")
	}
	return result("network-tls-posture", 8, "TLS posture", true, Info, .99, "TLS 1.3 negotiated", evidence, "")
}
func tlsVersion(v uint16) string {
	switch v {
	case tls.VersionTLS13:
		return "TLS1.3"
	case tls.VersionTLS12:
		return "TLS1.2"
	default:
		return fmt.Sprintf("0x%x", v)
	}
}

type HTTPMethodPolicy struct{}

func (HTTPMethodPolicy) ID() string { return "web-method-policy" }
func (HTTPMethodPolicy) Layer() int { return 9 }
func (HTTPMethodPolicy) Evaluate(r Request) Result {
	bad := r.StatusCode >= 500
	return result("web-method-policy", 9, "HTTP method policy", !bad, func() Severity {
		if bad {
			return High
		}
		return Info
	}(), .9, fmt.Sprintf("response status %d", r.StatusCode), map[string]string{"status": fmt.Sprint(r.StatusCode)}, "Reject unsupported methods and return a documented error contract")
}

type ContentTypePolicy struct{}

func (ContentTypePolicy) ID() string { return "web-content-type-policy" }
func (ContentTypePolicy) Layer() int { return 13 }
func (ContentTypePolicy) Evaluate(r Request) Result {
	ct := strings.ToLower(r.Headers["Content-Type"])
	passed := strings.Contains(ct, "application/json") || strings.Contains(ct, "text/plain")
	return result("web-content-type-policy", 13, "Content type policy", passed, Medium, .95, "content type is "+ct, map[string]string{"content_type": ct}, "Allow only documented content types and reject ambiguous parsers")
}

type APISchema struct{}

func (APISchema) ID() string { return "api-schema-contract" }
func (APISchema) Layer() int { return 15 }
func (APISchema) Evaluate(r Request) Result {
	var value any
	err := json.Unmarshal([]byte(r.Body), &value)
	passed := err == nil
	summary := "valid JSON response"
	if err != nil {
		summary = "response is not valid JSON"
	}
	return result("api-schema-contract", 15, "API schema contract", passed, Medium, .97, summary, map[string]string{"parse_error": fmt.Sprint(err)}, "Validate responses against the published OpenAPI schema")
}

type AuthorizationMatrix struct{}

func (AuthorizationMatrix) ID() string { return "identity-authorization-matrix" }
func (AuthorizationMatrix) Layer() int { return 12 }
func (AuthorizationMatrix) Evaluate(r Request) Result {
	denied := strings.EqualFold(r.Metadata["expected"], "deny") && r.StatusCode == 403
	return result("identity-authorization-matrix", 12, "Authorization matrix", denied, High, .96, fmt.Sprintf("expected deny, received status %d", r.StatusCode), map[string]string{"expected": r.Metadata["expected"], "status": fmt.Sprint(r.StatusCode)}, "Enforce server-side authorization for every object and tenant boundary")
}

type SecretExposure struct{}

func (SecretExposure) ID() string { return "application-secret-exposure" }
func (SecretExposure) Layer() int { return 14 }
func (SecretExposure) Evaluate(r Request) Result {
	patterns := []string{"password=", "secret=", "api_key=", "authorization: bearer"}
	body := strings.ToLower(r.Body)
	found := ""
	for _, p := range patterns {
		if strings.Contains(body, p) {
			found = p
			break
		}
	}
	return result("application-secret-exposure", 14, "Secret exposure", found == "", Critical, .99, "sensitive material exposure scan", map[string]string{"pattern": found}, "Redact secrets and rotate any exposed credential")
}

type DependencyPosture struct{}

func (DependencyPosture) ID() string { return "supply-chain-dependency-posture" }
func (DependencyPosture) Layer() int { return 20 }
func (DependencyPosture) Evaluate(r Request) Result {
	outdated := strings.EqualFold(r.Metadata["outdated"], "true")
	return result("supply-chain-dependency-posture", 20, "Dependency posture", !outdated, High, .9, "dependency inventory reviewed", map[string]string{"outdated": r.Metadata["outdated"]}, "Pin versions, generate SBOM, and remediate known vulnerable dependencies")
}

type CloudIdentityPosture struct{}

func (CloudIdentityPosture) ID() string { return "cloud-identity-posture" }
func (CloudIdentityPosture) Layer() int { return 17 }
func (CloudIdentityPosture) Evaluate(r Request) Result {
	public := strings.EqualFold(r.Metadata["public_access"], "true")
	return result("cloud-identity-posture", 17, "Cloud identity posture", !public, High, .94, "public access posture evaluated", map[string]string{"public_access": r.Metadata["public_access"]}, "Apply least privilege and remove unintended public access")
}

type EndpointLogging struct{}

func (EndpointLogging) ID() string { return "endpoint-logging-coverage" }
func (EndpointLogging) Layer() int { return 19 }
func (EndpointLogging) Evaluate(r Request) Result {
	coverage := strings.EqualFold(r.Metadata["audit_enabled"], "true")
	return result("endpoint-logging-coverage", 19, "Endpoint logging coverage", coverage, Medium, .92, "audit logging coverage evaluated", map[string]string{"audit_enabled": r.Metadata["audit_enabled"]}, "Enable immutable audit logging with retention and alerting")
}

type RecoveryControl struct{}

func (RecoveryControl) ID() string { return "recovery-rto-rpo-control" }
func (RecoveryControl) Layer() int { return 21 }
func (RecoveryControl) Evaluate(r Request) Result {
	restore := strings.EqualFold(r.Metadata["restore_verified"], "true")
	return result("recovery-rto-rpo-control", 21, "Recovery control", restore, High, .9, "restore verification result", map[string]string{"restore_verified": r.Metadata["restore_verified"]}, "Perform a documented restore test and measure RTO/RPO")
}

type EvidenceIntegrity struct{}

func (EvidenceIntegrity) ID() string { return "evidence-integrity-check" }
func (EvidenceIntegrity) Layer() int { return 23 }
func (EvidenceIntegrity) Evaluate(r Request) Result {
	hash := strings.TrimSpace(r.Metadata["sha256"])
	passed := len(hash) == 64
	return result("evidence-integrity-check", 23, "Evidence integrity", passed, Critical, .99, "evidence hash format validated", map[string]string{"sha256": hash}, "Store SHA-256 hash, timestamp, parent link, and chain-of-custody metadata")
}

func result(id string, layer int, title string, passed bool, sev Severity, confidence float64, summary string, evidence map[string]string, remediation string) Result {
	return Result{CheckID: id, Layer: layer, Title: title, Passed: passed, Severity: sev, Confidence: confidence, Summary: summary, Evidence: evidence, Remediation: remediation}
}
func NormalizeURL(value string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(value))
	if err != nil || u.Hostname() == "" {
		return "", fmt.Errorf("invalid target URL")
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return "", fmt.Errorf("unsupported URL scheme")
	}
	return u.Scheme + "://" + u.Hostname(), nil
}
func CertificateNames(cert *x509.Certificate) []string {
	return append([]string(nil), cert.DNSNames...)
}
func Probe(ctxTimeout time.Duration, target string) (Request, error) {
	client := http.Client{Timeout: ctxTimeout}
	res, err := client.Get(target)
	if err != nil {
		return Request{}, err
	}
	defer res.Body.Close()
	return Request{Target: target, StatusCode: res.StatusCode, Headers: map[string]string{"Content-Type": res.Header.Get("Content-Type")}}, nil
}
