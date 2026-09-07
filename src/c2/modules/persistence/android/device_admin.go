//go:build android

package persistence

import (
"os/exec"
)

func DeviceAdmin() bool {
cmd := exec.Command("dpm", "set-device-owner", "com.angel/.AdminReceiver")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func DeviceAdminRemove() bool {
cmd := exec.Command("dpm", "remove-active-admin", "com.angel/.AdminReceiver")
err := cmd.Run()
if err != nil {
return false
}
return true
}
