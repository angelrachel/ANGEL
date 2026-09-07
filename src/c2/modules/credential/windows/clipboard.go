//go:build windows

package windows

import (
"golang.org/x/sys/windows"
"unsafe"
)

type Clipboard struct{}

func NewClipboard() *Clipboard {
return &Clipboard{}
}

func (c *Clipboard) Capture() (string, error) {
user32 := windows.NewLazyDLL("user32.dll")
procOpenClipboard := user32.NewProc("OpenClipboard")
procGetClipboardData := user32.NewProc("GetClipboardData")
procCloseClipboard := user32.NewProc("CloseClipboard")

ret, _, _ := procOpenClipboard.Call(0)
if ret == 0 {
return "", nil
}
defer procCloseClipboard.Call()

handle, _, _ := procGetClipboardData.Call(1)
if handle == 0 {
return "", nil
}

dataPtr := (*uint16)(unsafe.Pointer(handle))
if dataPtr == nil {
return "", nil
}

data := windows.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(dataPtr))[:])
return data, nil
}
