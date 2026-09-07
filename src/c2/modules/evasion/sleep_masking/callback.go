package sleep_masking

import (
"syscall"
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

kernel32 := syscall.NewLazyDLL("kernel32.dll")
procCallNamedPipe := kernel32.NewProc("CallNamedPipeW")

procCallNamedPipe.Call(
uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("\\\\.\\pipe\\test"))),
0, 0, 0, 0, 0, 0,
)

<-done
return nil
}
