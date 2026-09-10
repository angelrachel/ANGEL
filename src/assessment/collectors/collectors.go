package collectors

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Observation struct {
	ID          string            `json:"id"`
	Kind        string            `json:"kind"`
	Target      string            `json:"target"`
	Attributes  map[string]string `json:"attributes"`
	Confidence  float64           `json:"confidence"`
	CollectedAt time.Time         `json:"collected_at"`
}
type HTTPCollector struct{ Client *http.Client }

func (c HTTPCollector) Collect(ctx context.Context, target string) (Observation, error) {
	u, err := url.Parse(target)
	if err != nil || u.Hostname() == "" || (u.Scheme != "https" && u.Scheme != "http") {
		return Observation{}, fmt.Errorf("invalid HTTP target")
	}
	client := c.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, target, nil)
	if err != nil {
		return Observation{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		return Observation{}, err
	}
	defer res.Body.Close()
	attrs := map[string]string{"status": fmt.Sprint(res.StatusCode), "content_type": res.Header.Get("Content-Type"), "server": res.Header.Get("Server"), "allow": res.Header.Get("Allow")}
	return Observation{ID: target + "/http", Kind: "http", Target: target, Attributes: attrs, Confidence: .98, CollectedAt: time.Now().UTC()}, nil
}

type TLSCollector struct{ Timeout time.Duration }

func (c TLSCollector) Collect(ctx context.Context, host string) (Observation, error) {
	timeout := c.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	dialer := tls.Dialer{Config: &tls.Config{MinVersion: tls.VersionTLS12, ServerName: host}}
	conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(host, "443"))
	if err != nil {
		return Observation{}, err
	}
	defer conn.Close()
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return Observation{}, fmt.Errorf("TLS collector received unexpected connection type")
	}
	state := tlsConn.ConnectionState()
	names := append([]string(nil), state.PeerCertificates[0].DNSNames...)
	sort.Strings(names)
	attrs := map[string]string{"version": fmt.Sprint(state.Version), "cipher": tls.CipherSuiteName(state.CipherSuite), "certificate_names": strings.Join(names, ",")}
	return Observation{ID: host + "/tls", Kind: "tls", Target: host, Attributes: attrs, Confidence: .99, CollectedAt: time.Now().UTC()}, nil
}

type DNSCollector struct{ Resolver *net.Resolver }

func (c DNSCollector) Collect(ctx context.Context, host string) ([]Observation, error) {
	resolver := c.Resolver
	if resolver == nil {
		resolver = net.DefaultResolver
	}
	ips, err := resolver.LookupHost(ctx, host)
	if err != nil {
		return nil, err
	}
	sort.Strings(ips)
	out := make([]Observation, 0, len(ips))
	for _, ip := range ips {
		out = append(out, Observation{ID: host + "/dns/" + ip, Kind: "dns", Target: host, Attributes: map[string]string{"ip": ip}, Confidence: .95, CollectedAt: time.Now().UTC()})
	}
	return out, nil
}

type Dependency struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Ecosystem string `json:"ecosystem"`
	Direct    bool   `json:"direct"`
}

func ParseDependencyManifest(data []byte) ([]Dependency, error) {
	var deps []Dependency
	if err := json.Unmarshal(data, &deps); err != nil {
		return nil, fmt.Errorf("invalid dependency manifest: %w", err)
	}
	for i := range deps {
		deps[i].Name = strings.TrimSpace(deps[i].Name)
		deps[i].Version = strings.TrimSpace(deps[i].Version)
		if deps[i].Name == "" || deps[i].Version == "" {
			return nil, fmt.Errorf("dependency %d is incomplete", i)
		}
	}
	sort.Slice(deps, func(i, j int) bool { return deps[i].Name < deps[j].Name })
	return deps, nil
}

type CloudPosture struct {
	PublicAccess      bool `json:"public_access"`
	MFARequired       bool `json:"mfa_required"`
	EncryptionEnabled bool `json:"encryption_enabled"`
	AuditEnabled      bool `json:"audit_enabled"`
}

func ParseCloudPosture(data []byte) (CloudPosture, error) {
	var posture CloudPosture
	if err := json.Unmarshal(data, &posture); err != nil {
		return posture, err
	}
	return posture, nil
}

type AuthorizationCase struct {
	Subject  string `json:"subject"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Expected string `json:"expected"`
	Actual   int    `json:"actual"`
}

func EvaluateAuthorizationMatrix(cases []AuthorizationCase) []Observation {
	sort.Slice(cases, func(i, j int) bool { return cases[i].Subject < cases[j].Subject })
	out := make([]Observation, 0, len(cases))
	for _, item := range cases {
		passed := (strings.EqualFold(item.Expected, "allow") && item.Actual < 400) || (strings.EqualFold(item.Expected, "deny") && item.Actual == 403)
		out = append(out, Observation{ID: item.Subject + ":" + item.Resource + ":" + item.Action, Kind: "authorization", Target: item.Resource, Attributes: map[string]string{"subject": item.Subject, "action": item.Action, "expected": item.Expected, "actual": fmt.Sprint(item.Actual), "passed": fmt.Sprint(passed)}, Confidence: .96, CollectedAt: time.Now().UTC()})
	}
	return out
}
