package task

import (
	"fmt"
	"sort"
	"time"
)

const (
	TaskQueued   = "QUEUED"
	TaskRunning  = "RUNNING"
	TaskComplete = "COMPLETE"
	TaskFailed   = "FAILED"
)

func (q *Queue) Enqueue(task *Task, now time.Time) error {
	if task == nil || task.ID == "" || task.Type == "" {
		return fmt.Errorf("task identity and type are required")
	}
	copy := cloneTask(task)
	copy.CreatedAt = copy.CreatedAt.UTC()
	if copy.CreatedAt.IsZero() {
		copy.CreatedAt = now.UTC()
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, exists := q.tasks[copy.ID]; exists {
		return fmt.Errorf("task already exists")
	}
	q.tasks[copy.ID] = copy
	return nil
}

type Claim struct {
	TaskID     string
	Attempts   int
	LeaseUntil time.Time
}

func (q *Queue) Claim(now time.Time, lease time.Duration) (*Task, Claim) {
	q.mu.Lock()
	defer q.mu.Unlock()
	var selected *Task
	for _, candidate := range q.tasks {
		if candidate == nil || candidate.Status == TaskComplete {
			continue
		}
		if candidate.Status == TaskQueued && candidate.LeaseUntil.After(now) {
			continue
		}
		if candidate.Status == TaskRunning && candidate.LeaseUntil.After(now) {
			continue
		}
		if selected == nil || candidate.CreatedAt.Before(selected.CreatedAt) {
			selected = candidate
		}
	}
	if selected == nil {
		return nil, Claim{}
	}
	selected.Status = TaskRunning
	selected.Attempts++
	selected.LeaseUntil = now.Add(lease).UTC()
	claimed := cloneTask(selected)
	return claimed, Claim{TaskID: selected.ID, Attempts: selected.Attempts, LeaseUntil: selected.LeaseUntil}
}

func (q *Queue) Ack(taskID string, now time.Time) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	item, ok := q.tasks[taskID]
	if !ok {
		return fmt.Errorf("task not found")
	}
	item.Status = TaskComplete
	item.LeaseUntil = time.Time{}
	_ = now
	return nil
}

func (q *Queue) Nack(taskID, reason string, retryAt time.Time) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	item, ok := q.tasks[taskID]
	if !ok {
		return fmt.Errorf("task not found")
	}
	item.Status = TaskQueued
	item.LastError = reason
	item.LeaseUntil = retryAt.UTC()
	return nil
}

func (q *Queue) Ready(now time.Time) []*Task {
	q.mu.RLock()
	defer q.mu.RUnlock()
	out := make([]*Task, 0)
	for _, item := range q.tasks {
		if item.Status == TaskQueued && !item.LeaseUntil.After(now) {
			out = append(out, cloneTask(item))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
