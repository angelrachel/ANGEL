package sleep_masking

import (
"syscall"
)

type StompingSleep struct{}

func NewStompingSleep() *StompingSleep {
return &StompingSleep{}
}

func (s *StompingSleep) Sleep(ms int) error {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procVirtualAlloc := kernel32.NewProc("VirtualAlloc")
procSleep := kernel32.NewProc("Sleep")

addr, _, _ := procVirtualAlloc.Call(0, 4096, 0x3000, 0x40)
procSleep.Call(uintptr(ms))
procVirtualAlloc.Call(addr, 4096, 0x8000, 0x40)

return nil
}
