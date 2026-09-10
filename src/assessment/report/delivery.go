package report

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"os"
)

type DeliveryResult struct {
	Path   string
	Status string
	Error  string
}

type Delivery struct {
	Key string
}

func (d Delivery) SaveJSON(path string, content string) DeliveryResult {
	return d.writeFile(path, []byte(content), 0644)
}

func (d Delivery) SaveMarkdown(path string, content string) DeliveryResult {
	return d.writeFile(path, []byte(content), 0644)
}

func (d Delivery) writeFile(path string, content []byte, mode os.FileMode) DeliveryResult {
	if path == "" {
		return DeliveryResult{Path: path, Status: "failed", Error: "output path is required"}
	}
	if err := os.WriteFile(path, content, mode); err != nil {
		return DeliveryResult{Path: path, Status: "failed", Error: err.Error()}
	}
	return DeliveryResult{Path: path, Status: "success"}
}

func (d Delivery) EncryptFile(path string, content []byte) DeliveryResult {
	if d.Key == "" {
		return DeliveryResult{Path: path, Status: "failed", Error: "encryption key is required"}
	}
	key := sha256.Sum256([]byte(d.Key))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return DeliveryResult{Path: path, Status: "failed", Error: err.Error()}
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return DeliveryResult{Path: path, Status: "failed", Error: err.Error()}
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return DeliveryResult{Path: path, Status: "failed", Error: err.Error()}
	}
	result := d.writeFile(path, gcm.Seal(nonce, nonce, content, nil), 0600)
	if result.Status != "success" && result.Error == "" {
		result.Error = errors.New("encrypted file write failed").Error()
	}
	return result
}
