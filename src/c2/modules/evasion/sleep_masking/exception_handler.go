package sleep_masking

import (
"syscall"
)

type ExceptionHandlerSleep struct{}

func NewExceptionHandlerSleep() *ExceptionHandlerSleep {
return &ExceptionHandlerSleep{}
}

func (e *ExceptionHandlerSleep) Sleep(ms int) error {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procAddVectoredExceptionHandler := kernel32.NewProc("AddVectoredExceptionHandler")
procSleep := kernel32.NewProc("Sleep")

handler := uintptr(1)
procAddVectoredExceptionHandler.Call(handler, uintptr(unsafe.Pointer(&exceptionHandler)))

procSleep.Call(uintptr(ms))
return nil
}

func exceptionHandler(exceptionInfo uintptr) uintptr {
return 1
}
