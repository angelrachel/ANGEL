//go:build windows

package rootkit

import (
"os/exec"
)

func SPIRead() bool {
cmd := exec.Command("cmd", "/c", "flashrom.exe -p internal -r firmware.bin")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func SPDump() bool {
cmd := exec.Command("cmd", "/c", "flashrom.exe -p internal -r bios.bin")
err := cmd.Run()
if err != nil {
return false
}
return true
}
