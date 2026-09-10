package task

import (
	"testing"
	"time"
)

func TestDurableQueueClaimAckAndRetry(t *testing.T) {
	now := time.Now().UTC()
	queue := NewQueue()
	if err := queue.Enqueue(&Task{ID: "task-1", Type: "assessment"}, now); err != nil {
		t.Fatal(err)
	}
	claimed, claim := queue.Claim(now, time.Minute)
	if claimed == nil || claim.Attempts != 1 || claimed.Status != TaskRunning {
		t.Fatalf("claim=%#v %#v", claimed, claim)
	}
	if ready := queue.Ready(now); len(ready) != 0 {
		t.Fatalf("leased task was ready: %#v", ready)
	}
	if err := queue.Nack("task-1", "worker timeout", now.Add(5*time.Minute)); err != nil {
		t.Fatal(err)
	}
	tooEarly, _ := queue.Claim(now, time.Minute)
	if tooEarly != nil {
		t.Fatal("backoff task claimed early")
	}
	retry, claim := queue.Claim(now.Add(5*time.Minute), time.Minute)
	if retry == nil || claim.Attempts != 2 {
		t.Fatalf("retry=%#v claim=%#v", retry, claim)
	}
	if err := queue.Ack("task-1", now); err != nil {
		t.Fatal(err)
	}
	if got := queue.Get("task-1"); got == nil || got.Status != TaskComplete {
		t.Fatalf("task=%#v", got)
	}
}

func TestDurableQueueRejectsDuplicateAndInvalidTask(t *testing.T) {
	queue := NewQueue()
	now := time.Now().UTC()
	if err := queue.Enqueue(nil, now); err == nil {
		t.Fatal("nil task accepted")
	}
	if err := queue.Enqueue(&Task{ID: "task-1"}, now); err == nil {
		t.Fatal("task without type accepted")
	}
	if err := queue.Enqueue(&Task{ID: "task-1", Type: "assessment"}, now); err != nil {
		t.Fatal(err)
	}
	if err := queue.Enqueue(&Task{ID: "task-1", Type: "assessment"}, now); err == nil {
		t.Fatal("duplicate task accepted")
	}
}
