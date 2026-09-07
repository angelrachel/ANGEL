//go:build windows

package syscall

import (
"golang.org/x/sys/windows"
"unsafe"
)

type IndirectSyscall struct{}

func NewIndirectSyscall() *IndirectSyscall {
return &IndirectSyscall{}
}

func (i *IndirectSyscall) Execute(syscallNumber uint16, args ...uintptr) (uintptr, error) {
var ret uintptr
var err error

kernel32 := windows.NewLazyDLL("kernel32.dll")
procGetProcAddress := kernel32.NewProc("GetProcAddress")
procGetModuleHandle := kernel32.NewProc("GetModuleHandleW")

moduleName, _ := windows.UTF16PtrFromString("ntdll.dll")
moduleHandle, _, _ := procGetModuleHandle.Call(uintptr(unsafe.Pointer(moduleName)))

funcName, _ := windows.UTF16PtrFromString("NtCreateProcessEx")
funcAddr, _, _ := procGetProcAddress.Call(moduleHandle, uintptr(unsafe.Pointer(funcName)))

ret, _, err = windows.Syscall6(
funcAddr,
uintptr(len(args)),
args[0], args[1], args[2], args[3], args[4], args[5],
)
return ret, err
}
