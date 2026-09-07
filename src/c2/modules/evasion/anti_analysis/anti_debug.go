//go:build windows

package anti_analysis

import (
"golang.org/x/sys/windows"
"unsafe"
)

type AntiDebug struct{}

func NewAntiDebug() *AntiDebug {
return &AntiDebug{}
}

func (a *AntiDebug) IsDebuggerPresent() bool {
kernel32 := windows.NewLazyDLL("kernel32.dll")
procIsDebuggerPresent := kernel32.NewProc("IsDebuggerPresent")
ret, _, _ := procIsDebuggerPresent.Call()
return ret != 0
}

func (a *AntiDebug) CheckPEB() bool {
kernel32 := windows.NewLazyDLL("kernel32.dll")
procGetCurrentProcess := kernel32.NewProc("GetCurrentProcess")
procNtQueryInformationProcess := kernel32.NewProc("NtQueryInformationProcess")

handle, _, _ := procGetCurrentProcess.Call()
var debugPort uintptr
procNtQueryInformationProcess.Call(handle, 7, uintptr(unsafe.Pointer(&debugPort)), 8, 0)
return debugPort != 0
}
