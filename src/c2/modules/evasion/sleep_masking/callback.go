//go:build windows

package sleep_masking

import (
"golang.org/x/sys/windows"
"time"
)

type CallbackSleep struct{}

func NewCallbackSleep() *CallbackSleep {
return &CallbackSleep{}
}

func (c *CallbackSleep) Sleep(ms int) error {
done := make(chan bool)
go func() {
time.Sleep(time.Duration(ms) * time.Millisecond)
done <- true
}()

kernel32 := windows.NewLazyDLL("kernel32.dll")
procCallNamedPipe := kernel32.NewProc("CallNamedPipeW")

name, _ := windows.UTF16PtrFromString("\\\\.\\pipe\\test")
procCallNamedPipe.Call(
uintptr(unsafe.Pointer(name)),
0, 0, 0, 0, 0, 0,
)

<-done
return nil
}
