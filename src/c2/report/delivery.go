package report

import (
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"crypto/sha256"
"io"
"os"
)

type DeliveryResult struct {
Path   string
Status string
}

type Delivery struct {
Key string
}

func (d Delivery) SaveJSON(path string, content string) DeliveryResult {
os.WriteFile(path, []byte(content), 0644)
return DeliveryResult{Path: path, Status: "success"}
}

func (d Delivery) SaveMarkdown(path string, content string) DeliveryResult {
os.WriteFile(path, []byte(content), 0644)
return DeliveryResult{Path: path, Status: "success"}
}

func (d Delivery) EncryptFile(path string, content []byte) DeliveryResult {
key := sha256.Sum256([]byte(d.Key))
block, _ := aes.NewCipher(key[:])
gcm, _ := cipher.NewGCM(block)
nonce := make([]byte, gcm.NonceSize())
io.ReadFull(rand.Reader, nonce)
os.WriteFile(path, gcm.Seal(nonce, nonce, content, nil), 0600)
return DeliveryResult{Path: path, Status: "success"}
}
