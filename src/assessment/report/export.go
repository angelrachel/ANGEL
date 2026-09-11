package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Export struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Findings  []Finding `json:"findings"`
}

func (r Report) JSON() (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}
	findings := append([]Finding(nil), r.Findings...)
	sort.SliceStable(findings, func(i, j int) bool {
		if strings.ToLower(findings[i].Severity) == strings.ToLower(findings[j].Severity) {
			return findings[i].Title < findings[j].Title
		}
		return severityRank(findings[i].Severity) > severityRank(findings[j].Severity)
	})
	payload, err := json.MarshalIndent(Export{ID: r.ID, Title: r.Title, StartTime: r.StartTime, EndTime: r.EndTime, Findings: findings}, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encode report: %w", err)
	}
	return string(payload), nil
}

func (r Report) ExecutiveSummary() string {
	counts := map[string]int{}
	for _, finding := range r.Findings {
		counts[strings.ToLower(strings.TrimSpace(finding.Severity))]++
	}
	return fmt.Sprintf("%d finding(s): %d critical, %d high, %d medium, %d low, %d info.", len(r.Findings), counts["critical"], counts["high"], counts["medium"], counts["low"], counts["info"])
}

func severityRank(value string) int {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "critical":
		return 5
	case "high":
		return 4
	case "medium":
		return 3
	case "low":
		return 2
	case "info":
		return 1
	default:
		return 0
	}
}
