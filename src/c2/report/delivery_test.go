package report

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeliveryReportsWriteErrors(t *testing.T) {
	result := (Delivery{}).SaveJSON("", "data")
	if result.Status != "failed" || result.Error == "" {
		t.Fatalf("expected path error, got %#v", result)
	}
	result = (Delivery{}).EncryptFile(filepath.Join(t.TempDir(), "out"), []byte("data"))
	if result.Status != "failed" || result.Error == "" {
		t.Fatalf("expected missing key error, got %#v", result)
	}
}

func TestDeliveryWritesEncryptedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.enc")
	result := (Delivery{Key: "test-key"}).EncryptFile(path, []byte("report"))
	if result.Status != "success" || result.Error != "" {
		t.Fatalf("expected encrypted write, got %#v", result)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat encrypted report: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected 0600 permissions, got %o", info.Mode().Perm())
	}
}
