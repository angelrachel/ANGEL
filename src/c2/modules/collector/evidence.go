//go:build windows

package collector

import (
"os/exec"
)

func CollectSystemInfo() string {
cmd := exec.Command("cmd", "/c", "systeminfo")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectUsers() string {
cmd := exec.Command("cmd", "/c", "net user")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectGroups() string {
cmd := exec.Command("cmd", "/c", "net localgroup")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectNetworkInfo() string {
cmd := exec.Command("cmd", "/c", "ipconfig /all")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectRunningProcesses() string {
cmd := exec.Command("cmd", "/c", "tasklist /v")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectServices() string {
cmd := exec.Command("cmd", "/c", "wmic service list brief")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectAll() string {
var result string
result += CollectSystemInfo()
result += CollectUsers()
result += CollectGroups()
result += CollectNetworkInfo()
result += CollectRunningProcesses()
result += CollectServices()
return result
}
