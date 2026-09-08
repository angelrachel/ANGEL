package crypto

import (
"crypto/hmac"
"crypto/sha256"
"encoding/hex"
)

type HMAC struct {
Key []byte
}

func NewHMAC(key string) *HMAC {
return &HMAC{Key: []byte(key)}
}

func (h *HMAC) Sign(data []byte) string {
mac := hmac.New(sha256.New, h.Key)
mac.Write(data)
return hex.EncodeToString(mac.Sum(nil))
}

func (h *HMAC) Verify(data []byte, signature string) bool {
expected := h.Sign(data)
return expected == signature
}
