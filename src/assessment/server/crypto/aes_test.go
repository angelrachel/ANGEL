package crypto

import "testing"

func TestAESGCMRoundTrip(t *testing.T) {
	cipher, err := NewAESGCM("test-key")
	if err != nil {
		t.Fatalf("new cipher failed: %v", err)
	}
	plaintext := []byte("authorized simulation metadata")
	encoded, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	decoded, err := cipher.Decrypt(encoded)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}
	if string(decoded) != string(plaintext) {
		t.Fatalf("round trip mismatch: %q", decoded)
	}
}

func TestAESGCMRejectsTruncatedAndTamperedCiphertext(t *testing.T) {
	cipher, err := NewAESGCM("test-key")
	if err != nil {
		t.Fatalf("new cipher failed: %v", err)
	}
	if _, err := cipher.Decrypt([]byte{1, 2}); err == nil {
		t.Fatal("expected truncated ciphertext error")
	}
	encoded, err := cipher.Encrypt([]byte("data"))
	if err != nil {
		t.Fatal(err)
	}
	encoded[len(encoded)-1] ^= 0xff
	if _, err := cipher.Decrypt(encoded); err == nil {
		t.Fatal("expected tampered ciphertext error")
	}
}
