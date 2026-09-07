package evidence

import (
"time"
)

type Timestamp struct {
Time time.Time
}

func NewTimestamp() *Timestamp {
return &Timestamp{
Time: time.Now().UTC(),
}
}

func (t *Timestamp) GetTimestamp() time.Time {
return t.Time
}

func (t *Timestamp) GetUnix() int64 {
return t.Time.Unix()
}

func (t *Timestamp) GetString() string {
return t.Time.Format(time.RFC3339)
}
