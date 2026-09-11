package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"time"
)

type AESGCM struct {
	AEAD cipher.AEAD
}

func NewAESGCM(key string) (*AESGCM, error) {
	hash := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, fmt.Errorf("aes cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("gcm: %w", err)
	}
	return &AESGCM{AEAD: gcm}, nil
}

func (a *AESGCM) Encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, a.AEAD.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return a.AEAD.Seal(nonce, nonce, data, nil), nil
}

func (a *AESGCM) Decrypt(data []byte) ([]byte, error) {
	ns := a.AEAD.NonceSize()
	if len(data) < ns {
		return nil, errors.New("ciphertext is shorter than nonce")
	}
	nonce, payload := data[:ns], data[ns:]
	return a.AEAD.Open([]byte{}, nonce, payload, nil)
}

func (a *AESGCM) EncryptWithTimestamp(data []byte) ([]byte, error) {
	ts := time.Now().Unix()
	payload := make([]byte, 8)
	for i := 0; i < 8; i++ {
		payload[i] = byte(ts >> (8 * i))
	}
	payload = append(payload, data...)
	return a.Encrypt(payload)
}
