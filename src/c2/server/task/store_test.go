package task

import "testing"

func TestQueueCopiesTasks(t *testing.T) {
	queue := NewQueue()
	original := &Task{ID: "task-1", Payload: map[string]interface{}{"mode": "simulation"}}
	queue.Add(original)
	original.Payload["mode"] = "changed"
	got := queue.Get("task-1")
	if got == nil || got.Payload["mode"] != "simulation" {
		t.Fatalf("expected queue to isolate input, got %#v", got)
	}
	got.Payload["mode"] = "mutated"
	if queue.Get("task-1").Payload["mode"] != "simulation" {
		t.Fatal("expected queue to isolate output")
	}
}

func TestResultStoreCopiesResultsAndRejectsNil(t *testing.T) {
	store := NewResultStore()
	store.Add(nil)
	store.Add(&Result{TaskID: "result-1", Output: map[string]interface{}{"status": "ok"}})
	got := store.Get("result-1")
	got.Output["status"] = "mutated"
	if store.Get("result-1").Output["status"] != "ok" {
		t.Fatal("expected result store to isolate output")
	}
}
