//go:build windows

package wmi

import (
"os/exec"
"strings"
)

func WMIQuery(targetHost, username, password, query string) string {
cmd := exec.Command("cmd", "/c", "wmic /node:"+targetHost+" /user:"+username+" /password:"+password+" "+query)
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func WMILogin(targetHost, username, password string) bool {
cmd := exec.Command("cmd", "/c", "wmic /node:"+targetHost+" /user:"+username+" /password:"+password+" os get caption")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func WMISystemInfo(targetHost, username, password string) string {
cmd := exec.Command("cmd", "/c", "wmic /node:"+targetHost+" /user:"+username+" /password:"+password+" os get caption,version,buildnumber")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}
