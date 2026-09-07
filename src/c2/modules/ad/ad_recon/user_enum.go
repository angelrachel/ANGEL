//go:build windows

package ad_recon

import (
"os/exec"
)

func EnumerateUser(username string) string {
cmd := exec.Command("cmd", "/c", "net user "+username+" /domain")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateUserSPN() string {
cmd := exec.Command("cmd", "/c", "setspn -Q */*")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateUserSessions() string {
cmd := exec.Command("cmd", "/c", "qwinsta")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateUserPrivileges() string {
cmd := exec.Command("cmd", "/c", "whoami /priv")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateUserGroups() string {
cmd := exec.Command("cmd", "/c", "whoami /groups")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func DumpUser(username string) string {
var result string
result += EnumerateUser(username)
result += EnumerateUserSPN()
result += EnumerateUserSessions()
result += EnumerateUserPrivileges()
result += EnumerateUserGroups()
return result
}
