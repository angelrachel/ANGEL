package report

import (
"time"
)

type RiskScore struct {
Score int
Level string
}

func NewRiskScore() *RiskScore {
return &RiskScore{}
}

func (r *RiskScore) Calculate(severity string, impact string) {
r.Score = len(severity) + len(impact)
if r.Score > 10 {
r.Level = "CRITICAL"
} else if r.Score > 5 {
r.Level = "HIGH"
} else {
r.Level = "MEDIUM"
}
}

func (r *RiskScore) GetScore() int {
return r.Score
}

func (r *RiskScore) GetLevel() string {
return r.Level
}

func (r *RiskScore) GetLastSeen() time.Time {
return time.Now()
}
