package reporting

type Metrics struct {
Severity      string
Confidence    string
Impact        string
BusinessImpact string
ROI           float64
}

func NewMetrics() *Metrics {
return &Metrics{}
}

func (m *Metrics) CalculateSeverity(exploitability, scope, privilege, dataImpact int) string {
score := exploitability + scope + privilege + dataImpact
if score >= 8 {
return "Critical (P0)"
} else if score >= 6 {
return "High (P1)"
} else if score >= 4 {
return "Medium (P2)"
} else if score >= 2 {
return "Low (P3)"
}
return "Info (P4)"
}

func (m *Metrics) CalculateConfidence(evidenceCount, reproducibility int) string {
score := evidenceCount + reproducibility
if score >= 8 {
return "High"
} else if score >= 5 {
return "Medium"
}
return "Low"
}

func (m *Metrics) CalculateROI(costToFix, costOfBreach int) float64 {
if costToFix == 0 {
return 0
}
return float64(costOfBreach-costToFix) / float64(costToFix) * 100
}
