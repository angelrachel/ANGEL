package report

import (
"time"
)

type Metrics struct {
Severity      string
Confidence    string
Impact        string
BusinessImpact string
ROI           string
}

func NewMetrics() *Metrics {
return &Metrics{}
}

func (m *Metrics) SetSeverity(severity string) {
m.Severity = severity
}

func (m *Metrics) SetConfidence(confidence string) {
m.Confidence = confidence
}

func (m *Metrics) SetImpact(impact string) {
m.Impact = impact
}

func (m *Metrics) SetBusinessImpact(businessImpact string) {
m.BusinessImpact = businessImpact
}

func (m *Metrics) SetROI(roi string) {
m.ROI = roi
}

func (m *Metrics) GetSeverity() string {
return m.Severity
}

func (m *Metrics) GetConfidence() string {
return m.Confidence
}

func (m *Metrics) GetImpact() string {
return m.Impact
}

func (m *Metrics) GetBusinessImpact() string {
return m.BusinessImpact
}

func (m *Metrics) GetROI() string {
return m.ROI
}

func (m *Metrics) UpdateLastSeen() {
_ = time.Now()
}
