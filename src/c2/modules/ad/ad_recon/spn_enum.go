//go:build windows

package ad_recon

import (
"os/exec"
)

func EnumerateSPNs() string {
cmd := exec.Command("cmd", "/c", "setspn -Q */*")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateSPNForUser(username string) string {
cmd := exec.Command("cmd", "/c", "setspn -L "+username)
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateSPNForService(service string) string {
cmd := exec.Command("cmd", "/c", "setspn -Q "+service+"/*")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func DumpAllSPNs() string {
return EnumerateSPNs()
}
