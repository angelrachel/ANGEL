package sleep_masking

import (
"syscall"
"unsafe"
)

type ThreadSpoof struct{}

func NewThreadSpoof() *ThreadSpoof {
return &ThreadSpoof{}
}

func (t *ThreadSpoof) Sleep(ms int) error {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procGetCurrentThread := kernel32.NewProc("GetCurrentThread")
procSetThreadStackGuarantee := kernel32.NewProc("SetThreadStackGuarantee")

currentThread, _, _ := procGetCurrentThread.Call()
guarantee := uintptr(4096)
procSetThreadStackGuarantee.Call(currentThread, guarantee)

syscall.Sleep(uint32(ms))
return nil
}
