//go:build windows

package sleep_masking

import (
"crypto/rc4"
"golang.org/x/sys/windows"
)

type EncryptFragmentsSleep struct{}

func NewEncryptFragmentsSleep() *EncryptFragmentsSleep {
return &EncryptFragmentsSleep{}
}

func (e *EncryptFragmentsSleep) Sleep(ms int) error {
key := []byte("angel-sleep-mask-key")
cipher, _ := rc4.NewCipher(key)

data := make([]byte, 4096)
cipher.XORKeyStream(data, data)

windows.Sleep(uint32(ms))

cipher.XORKeyStream(data, data)
return nil
}
