package report

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
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

// Validate checks the minimum structure required for a report to be shareable.
// It does not assert that a finding is factually correct; that remains a human
// review responsibility under the engagement process.
func (r Report) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.Title) == "" {
		return errors.New("report ID and title are required")
	}
	start, err := time.Parse(time.RFC3339, r.StartTime)
	if err != nil {
		return fmt.Errorf("invalid start time: %w", err)
	}
	end, err := time.Parse(time.RFC3339, r.EndTime)
	if err != nil {
		return fmt.Errorf("invalid end time: %w", err)
	}
	if end.Before(start) {
		return errors.New("report end time precedes start time")
	}
	for i, finding := range r.Findings {
		if strings.TrimSpace(finding.Title) == "" || strings.TrimSpace(finding.Description) == "" {
			return fmt.Errorf("finding %d requires title and description", i)
		}
		if !validSeverity(finding.Severity) {
			return fmt.Errorf("finding %d has unsupported severity %q", i, finding.Severity)
		}
	}
	return nil
}

func validSeverity(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "info", "low", "medium", "high", "critical":
		return true
	default:
		return false
	}
}

func escapeMarkdown(value string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\", "`", "\\`", "*", "\\*", "_", "\\_",
		"[", "\\[", "]", "\\]", "#", "\\#",
	)
	return replacer.Replace(value)
}

func (r *Report) Markdown() string {
	if r == nil {
		return ""
	}
	var output strings.Builder
	output.WriteString("# ")
	output.WriteString(escapeMarkdown(r.Title))
	output.WriteString("\n\n## Findings\n\n")
	for _, finding := range r.Findings {
		output.WriteString("### ")
		output.WriteString(escapeMarkdown(finding.Title))
		output.WriteString("\n**Severity:** ")
		output.WriteString(escapeMarkdown(finding.Severity))
		output.WriteString("\n**Description:** ")
		output.WriteString(escapeMarkdown(finding.Description))
		output.WriteString("\n**Evidence:** ")
		output.WriteString(escapeMarkdown(finding.Evidence))
		output.WriteString("\n\n")
	}
	return output.String()
}
