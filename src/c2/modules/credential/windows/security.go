package windows

import (
"os/exec"
)

type Security struct{}

func NewSecurity() *Security {
return &Security{}
}

func (s *Security) DumpSecurity() error {
cmd := exec.Command("reg", "save", "HKLM\\SECURITY", "security.hive")
return cmd.Run()
}
