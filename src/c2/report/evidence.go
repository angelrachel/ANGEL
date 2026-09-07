package report

import (
"time"
)

type Evidence struct {
ID        string
Timestamp time.Time
Data      string
}

func NewEvidence() *Evidence {
return &Evidence{
Timestamp: time.Now(),
}
}

func (e *Evidence) SetID(id string) {
e.ID = id
}

func (e *Evidence) SetData(data string) {
e.Data = data
}

func (e *Evidence) GetID() string {
return e.ID
}

func (e *Evidence) GetData() string {
return e.Data
}

func (e *Evidence) GetTimestamp() time.Time {
return e.Timestamp
}

func (e *Evidence) UpdateTimestamp() {
e.Timestamp = time.Now()
}
