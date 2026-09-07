package anti_analysis

import (
"syscall"
"unsafe"
)

type AntiDebug struct{}

func NewAntiDebug() *AntiDebug {
return &AntiDebug{}
}

func (a *AntiDebug) IsDebuggerPresent() bool {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procIsDebuggerPresent := kernel32.NewProc("IsDebuggerPresent")
ret, _, _ := procIsDebuggerPresent.Call()
return ret != 0
}

func (a *AntiDebug) CheckPEB() bool {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procGetCurrentProcess := kernel32.NewProc("GetCurrentProcess")
procNtQueryInformationProcess := kernel32.NewProc("NtQueryInformationProcess")

handle, _, _ := procGetCurrentProcess.Call()
var debugPort uintptr
procNtQueryInformationProcess.Call(handle, 7, uintptr(unsafe.Pointer(&debugPort)), 8, 0)

return debugPort != 0
}

func (a *AntiDebug) CheckNtGlobalFlag() bool {
pebAddr := uintptr(0x7ffdf000) // PEB address in x64
ntGlobalFlag := *(*uint32)(unsafe.Pointer(pebAddr + 0x68))
return ntGlobalFlag&0x70 != 0
}
