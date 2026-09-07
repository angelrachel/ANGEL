package crypto

import (
"crypto/aes"
"crypto/cipher"
"crypto/ecdh"
"crypto/rand"
"crypto/sha256"
"encoding/base64"
"errors"
"io"
)

// GenerateKeyPair generates ECDH P-256 key pair
func GenerateKeyPair() (*ecdh.PrivateKey, *ecdh.PublicKey, error) {
privateKey, err := ecdh.P256().GenerateKey(rand.Reader)
if err != nil {
return nil, nil, err
}
return privateKey, privateKey.PublicKey(), nil
}

// DeriveSharedKey derives shared key from private key and peer public key
func DeriveSharedKey(privateKey *ecdh.PrivateKey, peerPublicKey *ecdh.PublicKey) []byte {
sharedSecret, _ := privateKey.ECDH(peerPublicKey)
hash := sha256.Sum256(sharedSecret)
return hash[:]
}

// EncryptMessage encrypts plaintext with AES-256-GCM using shared key
func EncryptMessage(sharedKey []byte, plaintext []byte) (string, error) {
block, err := aes.NewCipher(sharedKey)
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

// DecryptMessage decrypts ciphertext with AES-256-GCM using shared key
func DecryptMessage(sharedKey []byte, ciphertextBase64 string) ([]byte, error) {
ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
if err != nil {
return nil, err
}

block, err := aes.NewCipher(sharedKey)
if err != nil {
return nil, err
}

gcm, err := cipher.NewGCM(block)
if err != nil {
return nil, err
}

nonceSize := gcm.NonceSize()
if len(ciphertext) < nonceSize {
return nil, errors.New("ciphertext too short")
}

nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
return gcm.Open(nil, nonce, ciphertext, nil)
}
