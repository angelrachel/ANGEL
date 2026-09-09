//go:build windows

package collector

import (
	"os/exec"
	"strings"
)

func MonitorProcesses() string {
	cmd := exec.Command("cmd", "/c", "tasklist /v")
	output, err := cmd.Output()
	if err != nil {
		return "Error: " + err.Error()
	}
	return strings.TrimSpace(string(output))
}

func MonitorNetwork() string {
	cmd := exec.Command("cmd", "/c", "netstat -ano")
	output, err := cmd.Output()
	if err != nil {
		return "Error: " + err.Error()
	}
	return strings.TrimSpace(string(output))
}

func MonitorConnections() string {
	cmd := exec.Command("cmd", "/c", "netstat -an")
	output, err := cmd.Output()
	if err != nil {
		return "Error: " + err.Error()
	}
	return strings.TrimSpace(string(output))
}
