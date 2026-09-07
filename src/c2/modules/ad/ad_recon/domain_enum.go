//go:build windows

package ad_recon

import (
"os/exec"
)

func EnumerateDomain() string {
cmd := exec.Command("cmd", "/c", "net view /domain")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateDomainController() string {
cmd := exec.Command("cmd", "/c", "nltest /dclist:angel.local")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateDomainUsers() string {
cmd := exec.Command("cmd", "/c", "net user /domain")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateDomainGroups() string {
cmd := exec.Command("cmd", "/c", "net group /domain")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateTrusts() string {
cmd := exec.Command("cmd", "/c", "nltest /domain_trusts")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func DumpAll() string {
var result string
result += EnumerateDomain()
result += EnumerateDomainController()
result += EnumerateDomainUsers()
result += EnumerateDomainGroups()
result += EnumerateTrusts()
return result
}
