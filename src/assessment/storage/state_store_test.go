package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
)

func TestFileStoreRoundTripAndDigest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "angel.json")
	store, empty, err := OpenFileStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if empty.Version != 1 {
		t.Fatalf("empty=%#v", empty)
	}
	state := PersistentState{Version: 1, Jobs: []domain.AssessmentJob{{ID: "job-1", Status: "CREATED", CreatedAt: time.Now().UTC()}}}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	reopened, loaded, err := OpenFileStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if reopened.Path() != path || len(loaded.Jobs) != 1 || loaded.Digest == "" {
		t.Fatalf("loaded=%#v", loaded)
	}
}

func TestFileStoreRejectsCorruptState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "angel.json")
	store, _, err := OpenFileStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(PersistentState{Version: 1}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"version":1,"digest":"bad"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := OpenFileStore(path); err == nil {
		t.Fatal("corrupt state accepted")
	}
}
