package query

import (
	"sort"
	"strings"

	"ANGEL/src/assessment/reporting"
)

type FindingFilter struct {
	Severity          string
	Status            string
	Target            string
	MinimumConfidence float64
	Search            string
	Page              int
	PageSize          int
}
type Page[T any] struct {
	Items    []T  `json:"items"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	Total    int  `json:"total"`
	HasNext  bool `json:"has_next"`
}

func Findings(items []reporting.Finding, filter FindingFilter) Page[reporting.Finding] {
	filtered := make([]reporting.Finding, 0, len(items))
	for _, item := range items {
		if filter.Severity != "" && !strings.EqualFold(item.Severity, filter.Severity) {
			continue
		}
		if filter.Status != "" && !strings.EqualFold(item.Status, filter.Status) {
			continue
		}
		if filter.Target != "" && !strings.EqualFold(item.Target, filter.Target) {
			continue
		}
		if item.Confidence < filter.MinimumConfidence {
			continue
		}
		if filter.Search != "" && !strings.Contains(strings.ToLower(item.Title+" "+item.Summary+" "+item.CheckID), strings.ToLower(filter.Search)) {
			continue
		}
		filtered = append(filtered, item)
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Severity == filtered[j].Severity {
			return filtered[i].ID < filtered[j].ID
		}
		return severityRank(filtered[i].Severity) > severityRank(filtered[j].Severity)
	})
	page := filter.Page
	if page < 1 {
		page = 1
	}
	size := filter.PageSize
	if size < 1 {
		size = 25
	}
	if size > 100 {
		size = 100
	}
	start := (page - 1) * size
	if start >= len(filtered) {
		return Page[reporting.Finding]{Items: []reporting.Finding{}, Page: page, PageSize: size, Total: len(filtered), HasNext: false}
	}
	end := start + size
	if end > len(filtered) {
		end = len(filtered)
	}
	return Page[reporting.Finding]{Items: filtered[start:end], Page: page, PageSize: size, Total: len(filtered), HasNext: end < len(filtered)}
}
func severityRank(value string) int {
	switch strings.ToUpper(value) {
	case "P0", "CRITICAL":
		return 5
	case "P1", "HIGH":
		return 4
	case "MEDIUM":
		return 3
	case "LOW":
		return 2
	default:
		return 1
	}
}
