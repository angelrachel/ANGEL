package plugins

import "sort"

type Module struct {
	ID         string `json:"id"`
	Layer      int    `json:"layer"`
	Name       string `json:"name"`
	Capability string `json:"capability"`
	Safe       bool   `json:"safe"`
}

var layerNames = []string{"Authorization & Engagement Control", "Identity, RBAC & Operator Security", "Scope Enforcement & Policy Engine", "Assessment Orchestration & Execution Control", "Platform Configuration & Secrets", "Asset Inventory", "Passive Reconnaissance", "Active Network Assessment", "Web Application Surface Mapping", "Cloud & External Exposure Mapping", "Authentication Security Assessment", "Authorization & Tenant Isolation Assessment", "Injection & Input Validation Assessment", "SSRF, File Handling & Server-Side Risk", "API, GraphQL & Business Logic Assessment", "Active Directory & Identity Posture", "Cloud IAM & Privilege Path", "Endpoint Security Posture & Execution Boundary", "Detection, Logging & Response Validation", "Vulnerability, Dependency & Supply Chain", "Controlled Validation Lab", "Safe Exploit Proof & Canary Validation", "Evidence Ledger & Chain of Custody", "Attack Path, Risk & P0/P1 Engine", "Reporting, Remediation & Retest"}
var capabilities = []string{"schema", "repository", "service", "validator", "parser", "normalizer", "matcher", "resolver", "planner", "dispatcher", "scheduler", "queue-adapter", "retry-policy", "timeout-guard", "cancellation", "capability-gate", "rate-limiter", "concurrency-guard", "state-machine", "event-logger", "evidence-collector", "redactor", "scorer", "correlation", "report-mapper", "replay-verifier"}

func Catalog() []Module {
	out := make([]Module, 0, 650)
	for layer := 1; layer <= 25; layer++ {
		for index, capability := range capabilities {
			out = append(out, Module{ID: id(layer, index+1), Layer: layer, Name: layerNames[layer-1] + " " + capability, Capability: capability, Safe: true})
		}
	}
	return out
}
func Layers() []string { return append([]string(nil), layerNames...) }
func Find(id string) (Module, bool) {
	for _, item := range Catalog() {
		if item.ID == id {
			return item, true
		}
	}
	return Module{}, false
}
func id(layer, index int) string { return "ANGEL-" + pad(layer) + "-" + pad(index) }
func pad(value int) string {
	if value < 10 {
		return "0" + string(rune('0'+value))
	}
	return string(rune('0'+value/10)) + string(rune('0'+value%10))
}
func Names() []string { out := append([]string(nil), layerNames...); sort.Strings(out); return out }
