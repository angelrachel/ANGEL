package task

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type State string

const (
	StateCreated     State = "created"
	StateApproved    State = "approved"
	StateQueued      State = "queued"
	StateRunning     State = "running"
	StateCompleted   State = "completed"
	StateFailed      State = "failed"
	StateRejected    State = "rejected"
	StateCancelled   State = "cancelled"
	StateExpired     State = "expired"
	StateCompensated State = "compensated"
)

type Lifecycle struct {
	TaskID    string    `json:"task_id"`
	State     State     `json:"state"`
	Version   int64     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
	Reason    string    `json:"reason,omitempty"`
}

func NewLifecycle(taskID string, now time.Time) (Lifecycle, error) {
	if strings.TrimSpace(taskID) == "" {
		return Lifecycle{}, errors.New("task ID is required")
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return Lifecycle{TaskID: taskID, State: StateCreated, Version: 1, UpdatedAt: now.UTC()}, nil
}

func (l Lifecycle) Transition(next State, reason string, now time.Time) (Lifecycle, error) {
	if strings.TrimSpace(l.TaskID) == "" {
		return Lifecycle{}, errors.New("task ID is required")
	}
	if l.Version < 1 {
		return Lifecycle{}, errors.New("lifecycle version must be positive")
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if !allowedTransition(l.State, next) {
		return Lifecycle{}, fmt.Errorf("invalid task transition from %q to %q", l.State, next)
	}
	if next == StateFailed || next == StateRejected || next == StateCancelled || next == StateExpired || next == StateCompensated {
		if strings.TrimSpace(reason) == "" {
			return Lifecycle{}, fmt.Errorf("reason is required for state %q", next)
		}
	}
	l.State = next
	l.Version++
	l.UpdatedAt = now.UTC()
	l.Reason = strings.TrimSpace(reason)
	return l, nil
}

func allowedTransition(from, to State) bool {
	switch from {
	case StateCreated:
		return to == StateApproved || to == StateRejected || to == StateExpired
	case StateApproved:
		return to == StateQueued || to == StateRejected || to == StateExpired || to == StateCancelled
	case StateQueued:
		return to == StateRunning || to == StateCancelled || to == StateExpired
	case StateRunning:
		return to == StateCompleted || to == StateFailed || to == StateCancelled || to == StateExpired
	case StateFailed, StateCancelled, StateExpired:
		return to == StateCompensated
	case StateCompleted:
		return false
	case StateRejected, StateCompensated:
		return false
	default:
		return false
	}
}
