//go:build windows

package credential

import (
"os/exec"
)

func CaptureClipboard() string {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Get-Clipboard\"")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CaptureClipboardLoop() string {
var result string
for i := 0; i < 10; i++ {
result += CaptureClipboard()
}
return result
}

func SetClipboard(data string) bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Set-Clipboard -Value '"+data+"'\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}
