package reporting

import (
"encoding/json"
"fmt"
"time"
)

type TechnicalReport struct {
Title          string
Date           string
Findings       []Finding
ExecutiveSummary string
Timeline       []TimelineEvent
}

type Finding struct {
ID          string
Title       string
Description string
Severity    string
Confidence  string
Impact      string
Evidence    []string
Reproduction string
Remediation string
}

type TimelineEvent struct {
Time    string
Event   string
Details string
}

func NewTechnicalReport(title string) *TechnicalReport {
return &TechnicalReport{
Title: title,
Date:  time.Now().Format("2006-01-02 15:04:05"),
}
}

func (t *TechnicalReport) AddFinding(finding Finding) {
t.Findings = append(t.Findings, finding)
}

func (t *TechnicalReport) AddTimelineEvent(event, details string) {
t.Timeline = append(t.Timeline, TimelineEvent{
Time:    time.Now().Format("15:04:05"),
Event:   event,
Details: details,
})
}

func (t *TechnicalReport) Generate() string {
jsonData, _ := json.MarshalIndent(t, "", "  ")
return string(jsonData)
}

func (t *TechnicalReport) GenerateMarkdown() string {
md := fmt.Sprintf("# %s\n\n", t.Title)
md += fmt.Sprintf("**Date:** %s\n\n", t.Date)
md += "## Executive Summary\n\n"
md += t.ExecutiveSummary + "\n\n"
md += "## Findings\n\n"

for _, f := range t.Findings {
md += fmt.Sprintf("### [%s] %s\n\n", f.Severity, f.Title)
md += fmt.Sprintf("**Description:** %s\n\n", f.Description)
md += fmt.Sprintf("**Impact:** %s\n\n", f.Impact)
md += fmt.Sprintf("**Remediation:** %s\n\n", f.Remediation)
}
return md
}
