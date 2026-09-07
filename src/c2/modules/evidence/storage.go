package evidence

import (
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"encoding/json"
"io"
"os"
)

type Storage struct {
Path     string
Key      []byte
}

func NewStorage(path string) *Storage {
key := make([]byte, 32)
io.ReadFull(rand.Reader, key)
return &Storage{
Path: path,
Key:  key,
}
}

func (s *Storage) StoreLocal(data interface{}) error {
jsonData, err := json.Marshal(data)
if err != nil {
return err
}
return os.WriteFile(s.Path, jsonData, 0644)
}

func (s *Storage) EncryptStore(data interface{}) error {
jsonData, err := json.Marshal(data)
if err != nil {
return err
}
block, err := aes.NewCipher(s.Key)
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
ciphertext := gcm.Seal(nonce, nonce, jsonData, nil)
return os.WriteFile(s.Path+".encrypted", ciphertext, 0644)
}

func (s *Storage) LoadLocal() ([]byte, error) {
return os.ReadFile(s.Path)
}

func (s *Storage) DecryptLoad() ([]byte, error) {
ciphertext, err := os.ReadFile(s.Path + ".encrypted")
if err != nil {
return nil, err
}
block, err := aes.NewCipher(s.Key)
if err != nil {
return nil, err
}
gcm, err := cipher.NewGCM(block)
if err != nil {
return nil, err
}
nonceSize := gcm.NonceSize()
nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
return gcm.Open(nil, nonce, ciphertext, nil)
}
