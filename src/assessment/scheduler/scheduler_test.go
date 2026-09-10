package scheduler

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

func TestSchedulerRetriesUntilSuccess(t *testing.T) {
	var attempts atomic.Int32
	scheduler := New(1, func(context.Context, Job) error {
		if attempts.Add(1) < 3 {
			return errors.New("temporary")
		}
		return nil
	})
	if err := scheduler.Enqueue(Job{ID: "job-1", MaxAttempts: 4}); err != nil {
		t.Fatal(err)
	}
	if err := scheduler.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 3 {
		t.Fatalf("attempts=%d", attempts.Load())
	}
}
func TestSchedulerDropsExhaustedJob(t *testing.T) {
	var attempts atomic.Int32
	scheduler := New(1, func(context.Context, Job) error { attempts.Add(1); return errors.New("permanent") })
	if err := scheduler.Enqueue(Job{ID: "job-1", MaxAttempts: 2}); err != nil {
		t.Fatal(err)
	}
	if err := scheduler.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 2 || scheduler.Pending() != 0 {
		t.Fatalf("attempts=%d pending=%d", attempts.Load(), scheduler.Pending())
	}
}
