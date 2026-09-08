//go:build windows

package collector

import (
"os/exec"
"strings"
)

func CollectNetworkInfo() string {
cmd := exec.Command("cmd", "/c", "ipconfig /all")
output, err := cmd.Output()
if err != nil {
return ""
}
return strings.TrimSpace(string(output))
}

func CollectInstalledSoftware() string {
cmd := exec.Command("cmd", "/c", "wmic product get name,version")
output, err := cmd.Output()
if err != nil {
return ""
}
return strings.TrimSpace(string(output))
}

func CollectAllEvidence() string {
var result strings.Builder
result.WriteString(CollectNetworkInfo())
result.WriteString(CollectInstalledSoftware())
return result.String()
}
