//go:build windows

package sleep_masking

import (
"golang.org/x/sys/windows"
"unsafe"
)

type ExceptionHandlerSleep struct{}

func NewExceptionHandlerSleep() *ExceptionHandlerSleep {
return &ExceptionHandlerSleep{}
}

func (e *ExceptionHandlerSleep) Sleep(ms int) error {
kernel32 := windows.NewLazyDLL("kernel32.dll")
procAddVectoredExceptionHandler := kernel32.NewProc("AddVectoredExceptionHandler")
procSleep := kernel32.NewProc("Sleep")

handler := uintptr(1)
// In real implementation, would add exception handler
procAddVectoredExceptionHandler.Call(handler, uintptr(unsafe.Pointer(&exceptionHandlerFunc)))

procSleep.Call(uintptr(ms))
return nil
}

func exceptionHandlerFunc(exceptionInfo uintptr) uintptr {
return 1
}
