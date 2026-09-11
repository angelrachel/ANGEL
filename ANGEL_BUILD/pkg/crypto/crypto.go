package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

type Session struct {
	Key []byte
}

func GenerateServerKey() (string, *ecdh.PrivateKey, error) {
	privateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return "", nil, err
	}
	pubB64 := base64.StdEncoding.EncodeToString(privateKey.PublicKey().Bytes())
	return pubB64, privateKey, nil
}

func ComputeShared(agentPubB64 string, privateKey *ecdh.PrivateKey) ([]byte, error) {
	agentPubRaw, err := base64.StdEncoding.DecodeString(agentPubB64)
	if err != nil {
		return nil, err
	}
	agentPub, err := ecdh.P256().NewPublicKey(agentPubRaw)
	if err != nil {
		return nil, err
	}
	shared, err := privateKey.ECDH(agentPub)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(shared)
	return hash[:], nil
}

func NewClientSession(serverPubB64 string) (*Session, string, error) {
	privateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, "", err
	}
	serverPubRaw, err := base64.StdEncoding.DecodeString(serverPubB64)
	if err != nil {
		return nil, "", err
	}
	serverPub, err := ecdh.P256().NewPublicKey(serverPubRaw)
	if err != nil {
		return nil, "", err
	}
	shared, err := privateKey.ECDH(serverPub)
	if err != nil {
		return nil, "", err
	}
	hash := sha256.Sum256(shared)
	ourPubB64 := base64.StdEncoding.EncodeToString(privateKey.PublicKey().Bytes())
	return &Session{Key: hash[:]}, ourPubB64, nil
}

func (s *Session) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *Session) Decrypt(ciphertextB64 string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
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
	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
