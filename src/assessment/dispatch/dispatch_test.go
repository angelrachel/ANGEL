package dispatch

import (
	"context"
	"testing"
	"time"
)

func TestBusPublishAndUnsubscribe(t *testing.T) {
	bus := NewBus()
	ch, closeSub := bus.Subscribe("JOB", 1)
	bus.Publish(Event{Type: "JOB", JobID: "job-1"})
	select {
	case event := <-ch:
		if event.JobID != "job-1" {
			t.Fatal(event)
		}
	case <-time.After(time.Second):
		t.Fatal("event timeout")
	}
	closeSub()
	bus.Publish(Event{Type: "JOB", JobID: "job-2"})
}
func TestPoolExecutesWork(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool := NewPool(2, 2)
	group := pool.Start(ctx)
	done := make(chan struct{})
	if err := pool.Submit(ctx, func(context.Context) error { close(done); return nil }); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("work timeout")
	}
	cancel()
	group.Wait()
}
