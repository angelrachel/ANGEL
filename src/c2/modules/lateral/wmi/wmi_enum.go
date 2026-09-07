//go:build windows

package wmi

import (
"os/exec"
"strings"
)

func WMIProcessEnum(targetHost, username, password string) string {
cmd := exec.Command("cmd", "/c", "wmic /node:"+targetHost+" /user:"+username+" /password:"+password+" process list brief")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func WMIServiceEnum(targetHost, username, password string) string {
cmd := exec.Command("cmd", "/c", "wmic /node:"+targetHost+" /user:"+username+" /password:"+password+" service list brief")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func WMINetworkConfig(targetHost, username, password string) string {
cmd := exec.Command("cmd", "/c", "wmic /node:"+targetHost+" /user:"+username+" /password:"+password+" nicconfig get ipaddress,macaddress")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}
