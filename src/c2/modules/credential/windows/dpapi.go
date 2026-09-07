//go:build windows

package credential

import (
"syscall"
"unsafe"
)

var (
crypt32                = syscall.NewLazyDLL("crypt32.dll")
procCryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
)

type DataBlob struct {
Size uint32
Data *byte
}

func DecryptDPAPI(data []byte) ([]byte, error) {
var input DataBlob
var output DataBlob
input.Data = &data[0]
input.Size = uint32(len(data))
ret, _, _ := procCryptUnprotectData.Call(
uintptr(unsafe.Pointer(&input)),
0,
0,
0,
0,
0,
uintptr(unsafe.Pointer(&output)),
)
if ret == 0 {
return nil, syscall.EINVAL
}
defer syscall.LocalFree(uintptr(unsafe.Pointer(output.Data)))
return unsafe.Slice(output.Data, output.Size), nil
}

func DecryptDPAPIWithEntropy(data, entropy []byte) ([]byte, error) {
var input DataBlob
var output DataBlob
var entropyBlob DataBlob
input.Data = &data[0]
input.Size = uint32(len(data))
entropyBlob.Data = &entropy[0]
entropyBlob.Size = uint32(len(entropy))
ret, _, _ := procCryptUnprotectData.Call(
uintptr(unsafe.Pointer(&input)),
0,
uintptr(unsafe.Pointer(&entropyBlob)),
0,
0,
0,
uintptr(unsafe.Pointer(&output)),
)
if ret == 0 {
return nil, syscall.EINVAL
}
defer syscall.LocalFree(uintptr(unsafe.Pointer(output.Data)))
return unsafe.Slice(output.Data, output.Size), nil
}
