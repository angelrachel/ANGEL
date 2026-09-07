//go:build android

package persistence

import (
"os/exec"
)

func MagiskModule() bool {
cmd := exec.Command("magisk", "--install-module", "angel.zip")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func MagiskModuleRemove() bool {
cmd := exec.Command("magisk", "--remove-modules")
err := cmd.Run()
if err != nil {
return false
}
return true
}
