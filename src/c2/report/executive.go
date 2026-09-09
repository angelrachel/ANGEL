package report

import "time"

type ExecutiveReport struct {
	Title          string
	Description    string
	RiskScore      int
	Impact         string
	Recommendation string
	Timestamp      string
}

func NewExecutiveReport() ExecutiveReport {
	return ExecutiveReport{Timestamp: time.Now().UTC().Format(time.RFC3339)}
}
