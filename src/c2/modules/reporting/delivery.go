package reporting

import (
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"encoding/json"
"io"
"os"
)

type Delivery struct {
Format string
Key    []byte
}

func NewDelivery(format string) *Delivery {
key := make([]byte, 32)
io.ReadFull(rand.Reader, key)
return &Delivery{
Format: format,
Key:    key,
}
}

func (d *Delivery) ExportJSON(data interface{}, filename string) error {
jsonData, err := json.MarshalIndent(data, "", "  ")
if err != nil {
return err
}
return os.WriteFile(filename+".json", jsonData, 0644)
}

func (d *Delivery) ExportMarkdown(data string, filename string) error {
return os.WriteFile(filename+".md", []byte(data), 0644)
}

func (d *Delivery) ExportPDF(data string, filename string) error {
// Placeholder - in real implementation would use PDF library
return os.WriteFile(filename+".pdf", []byte(data), 0644)
}

func (d *Delivery) EncryptExport(data []byte, filename string) error {
block, err := aes.NewCipher(d.Key)
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
return os.WriteFile(filename+".encrypted", ciphertext, 0644)
}
