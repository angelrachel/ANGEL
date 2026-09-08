package report

import (
"crypto/sha256"
"encoding/hex"
"time"
)

type Finding struct {
Title       string
Severity    string
Description string
Evidence    string
}

type Report struct {
ID        string
Title     string
StartTime string
EndTime   string
Findings  []Finding
}

func GenerateReport(title string, findings []Finding) Report {
now := time.Now().UTC()
hash := sha256.Sum256([]byte(title + now.String()))
return Report{
ID:        hex.EncodeToString(hash[:]),
Title:     title,
StartTime: now.Format(time.RFC3339),
EndTime:   now.Format(time.RFC3339),
Findings:  findings,
}
}

func (r *Report) Markdown() string {
output := "# " + r.Title + "\n\n"
output += "## Findings\n\n"
for _, finding := range r.Findings {
output += "### " + finding.Title + "\n"
output += "**Severity:** " + finding.Severity + "\n"
output += "**Description:** " + finding.Description + "\n"
output += "**Evidence:** " + finding.Evidence + "\n\n"
}
return output
}
