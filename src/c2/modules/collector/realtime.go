//go:build windows

package collector

import (
"os/exec"
)

func GetSystemTime() string {
cmd := exec.Command("cmd", "/c", "time /t")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func GetCurrentDirectory() string {
cmd := exec.Command("cmd", "/c", "cd")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func GetEnvironmentVariables() string {
cmd := exec.Command("cmd", "/c", "set")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func GetInstalledPrograms() string {
cmd := exec.Command("cmd", "/c", "wmic product get name,version")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func GetNetworkShares() string {
cmd := exec.Command("cmd", "/c", "net share")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func GetLoggedInUsers() string {
cmd := exec.Command("cmd", "/c", "query user")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func DumpAll() string {
var result string
result += GetSystemTime()
result += GetCurrentDirectory()
result += GetEnvironmentVariables()
result += GetInstalledPrograms()
result += GetNetworkShares()
result += GetLoggedInUsers()
return result
}
