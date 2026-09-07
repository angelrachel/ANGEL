//go:build windows

package ad_recon

import (
"os/exec"
"strings"
)

func EnumerateDomain() string {
cmd := exec.Command("cmd", "/c", "net view /domain")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func EnumerateDomainController() string {
cmd := exec.Command("cmd", "/c", "nltest /dclist:angel.local")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func EnumerateDomainUsers() string {
cmd := exec.Command("cmd", "/c", "net user /domain")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func EnumerateDomainGroups() string {
cmd := exec.Command("cmd", "/c", "net group /domain")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func EnumerateTrusts() string {
cmd := exec.Command("cmd", "/c", "nltest /domain_trusts")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}
