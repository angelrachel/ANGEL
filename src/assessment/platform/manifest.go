package platform

import (
	"fmt"
	"strings"
)

const Name = "ANGEL"
const Version = "1.0.0"

var layers = []string{
	"Authorization & Engagement Control",
	"Identity, RBAC & Operator Security",
	"Scope Enforcement & Policy Engine",
	"Assessment Orchestration & Execution Control",
	"Platform Configuration & Secrets",
	"Asset Inventory",
	"Passive Reconnaissance",
	"Active Network Assessment",
	"Web Application Surface Mapping",
	"Cloud & External Exposure Mapping",
	"Authentication Security Assessment",
	"Authorization & Tenant Isolation Assessment",
	"Injection & Input Validation Assessment",
	"SSRF, File Handling & Server-Side Risk",
	"API, GraphQL & Business Logic Assessment",
	"Active Directory & Identity Posture",
	"Cloud IAM & Privilege Path",
	"Endpoint Security Posture & Execution Boundary",
	"Detection, Logging & Response Validation",
	"Vulnerability, Dependency & Supply Chain",
	"Controlled Validation Lab",
	"Safe Exploit Proof & Canary Validation",
	"Evidence Ledger & Chain of Custody",
	"Attack Path, Risk & P0/P1 Engine",
	"Reporting, Remediation & Retest",
}

var deniedCapabilities = []string{
	"credential-collection", "credential-extraction", "persistence",
	"destructive-write", "log-deletion", "covert-channel",
	"process-injection", "evasion", "data-exfiltration", "arbitrary-command",
}

type Manifest struct {
	Name               string   `json:"name"`
	Version            string   `json:"version"`
	LayerCount         int      `json:"layer_count"`
	Layers             []string `json:"layers"`
	DeniedCapabilities []string `json:"denied_capabilities"`
	DefaultDecision    string   `json:"default_decision"`
}

func Current() Manifest {
	return Manifest{
		Name: Name, Version: Version, LayerCount: len(layers),
		Layers:             append([]string(nil), layers...),
		DeniedCapabilities: append([]string(nil), deniedCapabilities...),
		DefaultDecision:    "deny",
	}
}

func Layer(number int) (string, bool) {
	if number < 1 || number > len(layers) {
		return "", false
	}
	return layers[number-1], true
}

func IsDeniedCapability(capability string) bool {
	normalized := strings.ToLower(strings.TrimSpace(capability))
	for _, denied := range deniedCapabilities {
		if normalized == denied {
			return true
		}
	}
	return false
}

func (m Manifest) Validate() error {
	if m.Name != Name || m.Version == "" || m.LayerCount != 25 || len(m.Layers) != 25 {
		return fmt.Errorf("invalid ANGEL platform manifest")
	}
	if m.DefaultDecision != "deny" {
		return fmt.Errorf("platform default decision must be deny")
	}
	for _, capability := range deniedCapabilities {
		found := false
		for _, configured := range m.DeniedCapabilities {
			if configured == capability {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("denied capability missing: %s", capability)
		}
	}
	return nil
}
