package replay

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRepositoryFixturesReplay(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate replay package")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(sourceFile), "..", "..", "..", ".."))
	task, err := os.ReadFile(filepath.Join(root, "simulation", "fixtures", "tasks", "recon-http.json"))
	if err != nil {
		t.Fatalf("read task fixture: %v", err)
	}
	result, err := os.ReadFile(filepath.Join(root, "simulation", "fixtures", "reports", "recon-http-result.json"))
	if err != nil {
		t.Fatalf("read result fixture: %v", err)
	}
	if _, err := Replay(task, result); err != nil {
		t.Fatalf("replay repository fixtures: %v", err)
	}
}
