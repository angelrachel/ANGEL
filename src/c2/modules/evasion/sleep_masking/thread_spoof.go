//go:build windows

package sleep_masking

import (
"golang.org/x/sys/windows"
)

type ThreadSpoof struct{}

func NewThreadSpoof() *ThreadSpoof {
return &ThreadSpoof{}
}

func (t *ThreadSpoof) Sleep(ms int) error {
kernel32 := windows.NewLazyDLL("kernel32.dll")
procGetCurrentThread := kernel32.NewProc("GetCurrentThread")
procSetThreadStackGuarantee := kernel32.NewProc("SetThreadStackGuarantee")
procSleep := kernel32.NewProc("Sleep")

currentThread, _, _ := procGetCurrentThread.Call()
guarantee := uintptr(4096)
procSetThreadStackGuarantee.Call(currentThread, guarantee)

procSleep.Call(uintptr(ms))
return nil
}
