//go:build windows

package smb

import (
"os/exec"
"strings"
)

func SMBLogin(targetHost, username, password string) bool {
cmd := exec.Command("cmd", "/c", "net use \\\\"+targetHost+" /user:"+username+" "+password)
err := cmd.Run()
if err != nil {
return false
}
return true
}

func SMBLoginWithHash(targetHost, username, hash string) bool {
cmd := exec.Command("cmd", "/c", "net use \\\\"+targetHost+" /user:"+username+" "+hash)
err := cmd.Run()
if err != nil {
return false
}
return true
}

func SMBLogout(targetHost string) bool {
cmd := exec.Command("cmd", "/c", "net use \\\\"+targetHost+" /delete")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func SMBStatus(targetHost string) string {
cmd := exec.Command("cmd", "/c", "net use")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}
