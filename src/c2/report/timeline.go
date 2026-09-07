package report

import (
"time"
)

type Timeline struct {
StartTime time.Time
EndTime   time.Time
Events    []string
}

func NewTimeline() *Timeline {
return &Timeline{
StartTime: time.Now(),
Events:    []string{},
}
}

func (t *Timeline) AddEvent(event string) {
t.Events = append(t.Events, event)
}

func (t *Timeline) Close() {
t.EndTime = time.Now()
}

func (t *Timeline) GetEvents() []string {
return t.Events
}

func (t *Timeline) GetDuration() time.Duration {
return t.EndTime.Sub(t.StartTime)
}

func (t *Timeline) GetStartTime() time.Time {
return t.StartTime
}

func (t *Timeline) GetEndTime() time.Time {
return t.EndTime
}
