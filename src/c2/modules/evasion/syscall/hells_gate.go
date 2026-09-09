//go:build windows

package syscall

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

type HellsGate struct{}

func NewHellsGate() *HellsGate {
	return &HellsGate{}
}

func (h *HellsGate) GetSyscallNumber(functionName string) uint16 {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	procGetProcAddress := kernel32.NewProc("GetProcAddress")
	procGetModuleHandle := kernel32.NewProc("GetModuleHandleW")

	moduleName, _ := windows.UTF16PtrFromString("ntdll.dll")
	moduleHandle, _, _ := procGetModuleHandle.Call(uintptr(unsafe.Pointer(moduleName)))

	funcName, _ := windows.UTF16PtrFromString(functionName)
	funcAddr, _, _ := procGetProcAddress.Call(moduleHandle, uintptr(unsafe.Pointer(funcName)))

	if funcAddr == 0 {
		return 0
	}
	ssn := *(*uint16)(unsafe.Pointer(funcAddr + 4))
	return ssn
}
