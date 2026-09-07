package evidence

import (
"time"
)

type ReplayVerify struct {
Timestamp time.Time
Hash      string
}

func NewReplayVerify() *ReplayVerify {
return &ReplayVerify{
Timestamp: time.Now(),
}
}

func (r *ReplayVerify) Verify(hash string) bool {
return r.Hash == hash
}

func (r *ReplayVerify) SetHash(hash string) {
r.Hash = hash
}

func (r *ReplayVerify) GetTimestamp() time.Time {
return r.Timestamp
}
