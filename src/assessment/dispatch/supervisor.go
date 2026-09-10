package dispatch

import (
	"context"
	"fmt"
	"time"
)

type RetryPolicy struct {
	MaxAttempts int
	Timeout     time.Duration
	Backoff     time.Duration
}

type Execution struct {
	Attempts   int       `json:"attempts"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Err        error     `json:"-"`
}

// RunWithRetry executes one typed unit of work with bounded retries. The
// caller's context remains the source of truth for cancellation.
func RunWithRetry(ctx context.Context, policy RetryPolicy, work Work) Execution {
	execution := Execution{StartedAt: time.Now().UTC()}
	if work == nil {
		execution.Err = fmt.Errorf("work is required")
		execution.FinishedAt = time.Now().UTC()
		return execution
	}
	if policy.MaxAttempts < 1 {
		policy.MaxAttempts = 1
	}
	if policy.Backoff < 0 {
		policy.Backoff = 0
	}
	for execution.Attempts < policy.MaxAttempts {
		if err := ctx.Err(); err != nil {
			execution.Err = err
			break
		}
		execution.Attempts++
		attemptCtx := ctx
		cancel := func() {}
		if policy.Timeout > 0 {
			attemptCtx, cancel = context.WithTimeout(ctx, policy.Timeout)
		}
		err := work(attemptCtx)
		cancel()
		if err == nil {
			execution.Err = nil
			break
		}
		execution.Err = err
		if execution.Attempts == policy.MaxAttempts {
			break
		}
		timer := time.NewTimer(policy.Backoff)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			execution.Err = ctx.Err()
		case <-timer.C:
		}
	}
	execution.FinishedAt = time.Now().UTC()
	return execution
}
