package brain

type RiskAssessment struct {
RiskLevel   string
Score       int
Environment string
}

func NewRiskAssessment() *RiskAssessment {
return &RiskAssessment{
RiskLevel: "unknown",
Score:     50,
}
}

func (r *RiskAssessment) Assess(environment string, indicators []string) string {
score := 0
for _, indicator := range indicators {
switch indicator {
case "edr_present":
score += 30
case "av_present":
score += 10
case "sandbox":
score += 40
case "debugger":
score += 50
}
}

if score > 70 {
r.RiskLevel = "critical"
} else if score > 40 {
r.RiskLevel = "high"
} else if score > 20 {
r.RiskLevel = "medium"
} else {
r.RiskLevel = "low"
}
r.Score = score
return r.RiskLevel
}
