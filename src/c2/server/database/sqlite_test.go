package database

import (
	"sync"
	"testing"
	"time"
)

func TestListReturnsSnapshot(t *testing.T) {
	store := NewStore()
	store.Set("key", "value")
	snapshot := store.List()
	snapshot["key"] = "modified-outside-store"
	if got := store.Get("key"); got != "value" {
		t.Fatalf("expected store value to remain isolated, got %#v", got)
	}
}

func TestExpiryDoesNotDeleteReplacement(t *testing.T) {
	store := NewStore()
	store.SetWithExpiry("key", "old", 20*time.Millisecond)
	store.Set("key", "new")
	time.Sleep(40 * time.Millisecond)
	if got := store.Get("key"); got != "new" {
		t.Fatalf("expected replacement to survive old expiry timer, got %#v", got)
	}
}

func TestStoreConcurrentAccess(t *testing.T) {
	store := NewStore()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			store.Set("key", index)
			_ = store.Get("key")
			_ = store.List()
		}(i)
	}
	wg.Wait()
	if store.Get("key") == nil {
		t.Fatal("expected a value after concurrent writes")
	}
}
