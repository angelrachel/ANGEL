package report

import "time"

type ExecutiveSummary struct {
Title          string
Description    string
RiskScore      string
Impact         string
Recommendation string
}

func NewExecutiveSummary() *ExecutiveSummary {
return &ExecutiveSummary{}
}

func (e *ExecutiveSummary) SetTitle(title string) {
e.Title = title
}

func (e *ExecutiveSummary) SetDescription(description string) {
e.Description = description
}

func (e *ExecutiveSummary) SetRiskScore(riskScore string) {
e.RiskScore = riskScore
}

func (e *ExecutiveSummary) SetImpact(impact string) {
e.Impact = impact
}

func (e *ExecutiveSummary) SetRecommendation(recommendation string) {
e.Recommendation = recommendation
}

func (e *ExecutiveSummary) GetTitle() string {
return e.Title
}

func (e *ExecutiveSummary) GetDescription() string {
return e.Description
}

func (e *ExecutiveSummary) GetRiskScore() string {
return e.RiskScore
}

func (e *ExecutiveSummary) GetImpact() string {
return e.Impact
}

func (e *ExecutiveSummary) GetRecommendation() string {
return e.Recommendation
}

func (e *ExecutiveSummary) GetLastSeen() time.Time {
return time.Now()
}
