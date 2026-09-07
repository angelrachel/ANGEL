//go:build windows

package collector

import (
"os/exec"
)

func MonitorProcesses() string {
cmd := exec.Command("cmd", "/c", "tasklist /v")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func MonitorNetwork() string {
cmd := exec.Command("cmd", "/c", "netstat -ano")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func MonitorConnections() string {
cmd := exec.Command("cmd", "/c", "netstat -an")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func MonitorFileSystem() string {
cmd := exec.Command("cmd", "/c", "dir /s /b C:\\")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func MonitorAll() string {
var result string
result += MonitorProcesses()
result += MonitorNetwork()
result += MonitorConnections()
result += MonitorFileSystem()
return result
}
