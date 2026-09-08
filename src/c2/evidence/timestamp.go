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

func (t *Timestamp) GetTime() time.Time {
return t.Time
}

func (t *Timestamp) GetUnix() int64 {
return t.Time.Unix()
}

func (t *Timestamp) GetRFC3339() string {
return t.Time.Format(time.RFC3339)
}
