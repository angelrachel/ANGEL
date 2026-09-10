package dispatch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRunWithRetryEventuallySucceeds(t *testing.T) {
	attempts := 0
	result := RunWithRetry(context.Background(), RetryPolicy{MaxAttempts: 3, Backoff: time.Millisecond}, func(context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if result.Err != nil || result.Attempts != 3 {
		t.Fatalf("unexpected execution: %#v", result)
	}
}

func TestRunWithRetryHonorsTimeoutAndAttemptBound(t *testing.T) {
	result := RunWithRetry(context.Background(), RetryPolicy{MaxAttempts: 2, Timeout: 5 * time.Millisecond}, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if result.Err == nil || result.Attempts != 2 {
		t.Fatalf("timeout was not bounded: %#v", result)
	}
}

func TestRunWithRetryHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := RunWithRetry(ctx, RetryPolicy{MaxAttempts: 5}, func(context.Context) error { t.Fatal("cancelled work ran"); return nil })
	if !errors.Is(result.Err, context.Canceled) {
		t.Fatalf("err=%v", result.Err)
	}
	if result.Attempts != 0 {
		t.Fatalf("attempts=%d", result.Attempts)
	}
}
