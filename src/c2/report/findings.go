package report

import (
"time"
)

type Finding struct {
ID          string
Title       string
Description string
Severity    string
Evidence    string
Remediation string
Timestamp   time.Time
}

func NewFinding(id string) *Finding {
return &Finding{
ID:        id,
Timestamp: time.Now(),
}
}

func (f *Finding) SetTitle(title string) {
f.Title = title
}

func (f *Finding) SetDescription(description string) {
f.Description = description
}

func (f *Finding) SetSeverity(severity string) {
f.Severity = severity
}

func (f *Finding) SetEvidence(evidence string) {
f.Evidence = evidence
}

func (f *Finding) SetRemediation(remediation string) {
f.Remediation = remediation
}

func (f *Finding) GetTitle() string {
return f.Title
}

func (f *Finding) GetDescription() string {
return f.Description
}

func (f *Finding) GetSeverity() string {
return f.Severity
}

func (f *Finding) GetEvidence() string {
return f.Evidence
}

func (f *Finding) GetRemediation() string {
return f.Remediation
}

func (f *Finding) GetTimestamp() time.Time {
return f.Timestamp
}
