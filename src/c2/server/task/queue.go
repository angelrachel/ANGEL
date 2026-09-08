package task

import (
"sync"
"time"
)

type Task struct {
ID        string
AgentID   string
Type      string
Payload   map[string]interface{}
CreatedAt time.Time
}

type Queue struct {
mu    sync.Mutex
tasks map[string]*Task
}

func NewQueue() *Queue {
return &Queue{tasks: make(map[string]*Task)}
}

func (q *Queue) Add(task *Task) {
q.mu.Lock()
defer q.mu.Unlock()
q.tasks[task.ID] = task
}

func (q *Queue) Get(taskID string) *Task {
q.mu.Lock()
defer q.mu.Unlock()
return q.tasks[taskID]
}

func (q *Queue) List() []*Task {
q.mu.Lock()
defer q.mu.Unlock()
var tasks []*Task
for _, task := range q.tasks {
tasks = append(tasks, task)
}
return tasks
}

func (q *Queue) Delete(taskID string) {
q.mu.Lock()
defer q.mu.Unlock()
delete(q.tasks, taskID)
}
