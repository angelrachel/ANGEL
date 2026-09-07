package implant

import (
"crypto/rand"
"crypto/rsa"
"crypto/x509"
"encoding/base64"
"encoding/pem"
"os"
)

type ImplantGenerator struct {
ServerURL   string
CallbackInt int
Key         []byte
}

func NewImplantGenerator(serverURL string) *ImplantGenerator {
key := make([]byte, 32)
rand.Read(key)
return &ImplantGenerator{
ServerURL:   serverURL,
CallbackInt: 60,
Key:         key,
}
}

func (i *ImplantGenerator) SetCallbackInterval(seconds int) {
i.CallbackInt = seconds
}

func (i *ImplantGenerator) GenerateEncryptionKey() []byte {
key := make([]byte, 32)
rand.Read(key)
i.Key = key
return key
}

func (i *ImplantGenerator) EncryptPayload(payload []byte) string {
encrypted := make([]byte, len(payload))
for idx, b := range payload {
encrypted[idx] = b ^ i.Key[idx%len(i.Key)]
}
return base64.StdEncoding.EncodeToString(encrypted)
}

func (i *ImplantGenerator) DecryptPayload(encoded string) ([]byte, error) {
data, err := base64.StdEncoding.DecodeString(encoded)
if err != nil {
return nil, err
}
decrypted := make([]byte, len(data))
for idx, b := range data {
decrypted[idx] = b ^ i.Key[idx%len(i.Key)]
}
return decrypted, nil
}

func (i *ImplantGenerator) GenerateRSAPrivateKey() ([]byte, error) {
privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
if err != nil {
return nil, err
}
derBytes := x509.MarshalPKCS1PrivateKey(privateKey)
block := &pem.Block{
Type:  "RSA PRIVATE KEY",
Bytes: derBytes,
}
return pem.EncodeToMemory(block), nil
}

func (i *ImplantGenerator) SavePrivateKey(path string) bool {
keyBytes, err := i.GenerateRSAPrivateKey()
if err != nil {
return false
}
file, err := os.Create(path)
if err != nil {
return false
}
defer file.Close()
file.Write(keyBytes)
return true
}

func (i *ImplantGenerator) GetServerURL() string {
return i.ServerURL
}

func (i *ImplantGenerator) GetCallbackInterval() int {
return i.CallbackInt
}

func (i *ImplantGenerator) GetKey() []byte {
return i.Key
}
