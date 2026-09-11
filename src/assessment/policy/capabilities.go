package policy

import "sort"

var allowedCapabilities = []string{"surface-map", "tls-assessment", "api-contract", "authorization-matrix", "evidence-collection", "dependency-inventory", "detection-validation", "synthetic-canary", "lab-proof", "lab-agent-simulation"}
var deniedCapabilities = []string{"credential-collection", "credential-extraction", "persistence", "destructive-write", "log-deletion", "covert-channel", "process-injection", "evasion", "data-exfiltration", "arbitrary-command"}

// CapabilityPolicy returns copies so callers cannot mutate the process-wide policy.
func CapabilityPolicy() (allowed, denied []string) {
	allowed = append([]string(nil), allowedCapabilities...)
	denied = append([]string(nil), deniedCapabilities...)
	sort.Strings(allowed)
	sort.Strings(denied)
	return allowed, denied
}

func CapabilityAllowed(name string) bool {
	for _, candidate := range allowedCapabilities {
		if candidate == name {
			return true
		}
	}
	return false
}

func CapabilityDenied(name string) bool {
	for _, candidate := range deniedCapabilities {
		if candidate == name {
			return true
		}
	}
	return false
}
