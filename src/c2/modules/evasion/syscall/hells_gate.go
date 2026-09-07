package syscall

import (
"syscall"
"unsafe"
)

type HellsGate struct{}

func NewHellsGate() *HellsGate {
return &HellsGate{}
}

func (h *HellsGate) GetSyscallNumber(functionName string) uint16 {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procGetProcAddress := kernel32.NewProc("GetProcAddress")
procGetModuleHandle := kernel32.NewProc("GetModuleHandleW")

moduleName, _ := syscall.UTF16PtrFromString("ntdll.dll")
moduleHandle, _, _ := procGetModuleHandle.Call(uintptr(unsafe.Pointer(moduleName)))

funcName, _ := syscall.UTF16PtrFromString(functionName)
funcAddr, _, _ := procGetProcAddress.Call(moduleHandle, uintptr(unsafe.Pointer(funcName)))

if funcAddr == 0 {
return 0
}

// Extract SSN from function (first 4 bytes contain syscall number)
ssn := *(*uint16)(unsafe.Pointer(funcAddr + 4))
return ssn
}
