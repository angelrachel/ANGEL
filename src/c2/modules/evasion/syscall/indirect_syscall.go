package syscall

import (
"syscall"
"unsafe"
)

type IndirectSyscall struct{}

func NewIndirectSyscall() *IndirectSyscall {
return &IndirectSyscall{}
}

func (i *IndirectSyscall) Execute(syscallNumber uint16, args ...uintptr) (uintptr, error) {
var ret uintptr
var err error

// Get ntdll!Nt* function address
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procGetProcAddress := kernel32.NewProc("GetProcAddress")
procGetModuleHandle := kernel32.NewProc("GetModuleHandleW")

moduleName, _ := syscall.UTF16PtrFromString("ntdll.dll")
moduleHandle, _, _ := procGetModuleHandle.Call(uintptr(unsafe.Pointer(moduleName)))

funcName, _ := syscall.UTF16PtrFromString("NtCreateProcessEx")
funcAddr, _, _ := procGetProcAddress.Call(moduleHandle, uintptr(unsafe.Pointer(funcName)))

// Direct syscall using the address and SSN
ret, _, err = syscall.Syscall6(
funcAddr,
uintptr(len(args)),
args[0], args[1], args[2], args[3], args[4], args[5],
)

return ret, err
}
