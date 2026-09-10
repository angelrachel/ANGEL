package task

import (
	"sync"
	"time"
)

type Task struct {
	ID         string
	AgentID    string
	Type       string
	Payload    map[string]interface{}
	CreatedAt  time.Time
	Status     string
	Attempts   int
	LeaseUntil time.Time
	LastError  string
}

type Queue struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

func NewQueue() *Queue {
	return &Queue{tasks: make(map[string]*Task)}
}

func cloneTask(task *Task) *Task {
	if task == nil {
		return nil
	}
	clone := *task
	if task.Payload != nil {
		clone.Payload = make(map[string]interface{}, len(task.Payload))
		for key, value := range task.Payload {
			clone.Payload[key] = value
		}
	}
	return &clone
}

func (q *Queue) Add(task *Task) {
	if task == nil || task.ID == "" {
		return
	}
	copy := cloneTask(task)
	if copy.Status == "" {
		copy.Status = TaskQueued
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tasks[task.ID] = copy
}

func (q *Queue) Get(taskID string) *Task {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return cloneTask(q.tasks[taskID])
}

func (q *Queue) List() []*Task {
	q.mu.RLock()
	defer q.mu.RUnlock()
	tasks := make([]*Task, 0, len(q.tasks))
	for _, task := range q.tasks {
		tasks = append(tasks, cloneTask(task))
	}
	return tasks
}

func (q *Queue) Delete(taskID string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.tasks, taskID)
}
