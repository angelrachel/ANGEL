package brain

type RiskLevel struct {
Score int
Level string
}

type RiskAssessor struct {
Threshold int
}

func NewRiskAssessor(threshold int) *RiskAssessor {
return &RiskAssessor{Threshold: threshold}
}

func (r *RiskAssessor) Assess(score int) RiskLevel {
if score >= r.Threshold {
return RiskLevel{Score: score, Level: "high"}
}
return RiskLevel{Score: score, Level: "low"}
}

func (r *RiskAssessor) SetThreshold(threshold int) {
r.Threshold = threshold
}
