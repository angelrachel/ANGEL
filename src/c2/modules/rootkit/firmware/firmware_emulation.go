//go:build windows

package rootkit

import (
"os/exec"
)

func FirmwareEmulate() bool {
cmd := exec.Command("cmd", "/c", "qemu-system-x86_64 -m 512 -bios firmware.bin")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func FirmwareDump() bool {
cmd := exec.Command("cmd", "/c", "flashrom.exe -p internal -r dump.bin")
err := cmd.Run()
if err != nil {
return false
}
return true
}
