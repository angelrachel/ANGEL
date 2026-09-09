package brain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

type SleepMasker struct {
	AEAD cipher.AEAD
}

func NewSleepMasker(key string) *SleepMasker {
	hash := sha256.Sum256([]byte(key))
	block, _ := aes.NewCipher(hash[:])
	gcm, _ := cipher.NewGCM(block)
	return &SleepMasker{AEAD: gcm}
}

func (s *SleepMasker) EncryptRegion(data []byte) ([]byte, error) {
	nonce := make([]byte, s.AEAD.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	return s.AEAD.Seal(nonce, nonce, data, nil), nil
}

func (s *SleepMasker) DecryptRegion(data []byte) ([]byte, error) {
	ns := s.AEAD.NonceSize()
	if len(data) < ns {
		return []byte{}, nil
	}
	nonce, payload := data[:ns], data[ns:]
	return s.AEAD.Open([]byte{}, nonce, payload, nil)
}
