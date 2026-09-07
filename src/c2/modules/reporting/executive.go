package reporting

import (
"encoding/json"
"time"
)

type ExecutiveReport struct {
Title       string
Date        string
Summary     string
Impact      string
Recommendations []string
RiskScore   int
}

func NewExecutiveReport(title string) *ExecutiveReport {
return &ExecutiveReport{
Title: title,
Date:  time.Now().Format("2006-01-02 15:04:05"),
}
}

func (e *ExecutiveReport) AddRecommendation(rec string) {
e.Recommendations = append(e.Recommendations, rec)
}

func (e *ExecutiveReport) Generate() string {
jsonData, _ := json.MarshalIndent(e, "", "  ")
return string(jsonData)
}

func (e *ExecutiveReport) GenerateMarkdown() string {
md := "# " + e.Title + "\n\n"
md += "**Date:** " + e.Date + "\n\n"
md += "## Executive Summary\n\n" + e.Summary + "\n\n"
md += "## Impact\n\n" + e.Impact + "\n\n"
md += "## Risk Score: " + string(rune(e.RiskScore)) + "/100\n\n"
md += "## Recommendations\n\n"
for _, rec := range e.Recommendations {
md += "- " + rec + "\n"
}
return md
}
