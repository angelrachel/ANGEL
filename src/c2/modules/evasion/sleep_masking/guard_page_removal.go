//go:build windows

package sleep_masking

import (
"golang.org/x/sys/windows"
"unsafe"
)

type GuardPageRemoval struct{}

func NewGuardPageRemoval() *GuardPageRemoval {
return &GuardPageRemoval{}
}

func (g *GuardPageRemoval) Sleep(ms int) error {
kernel32 := windows.NewLazyDLL("kernel32.dll")
procVirtualProtect := kernel32.NewProc("VirtualProtect")
procSleep := kernel32.NewProc("Sleep")

var oldProtect uint32
addr := uintptr(0x1000)
procVirtualProtect.Call(addr, 4096, 0x40, uintptr(unsafe.Pointer(&oldProtect)))
procSleep.Call(uintptr(ms))
procVirtualProtect.Call(addr, 4096, oldProtect, uintptr(unsafe.Pointer(&oldProtect)))
return nil
}
