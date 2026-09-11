package query

import (
	"sort"
	"strings"

	"ANGEL/src/assessment/reporting"
)

type ReportFilter struct {
	EngagementID string
	Page         int
	PageSize     int
}

func Reports(items []reporting.Report, filter ReportFilter) Page[reporting.Report] {
	filtered := make([]reporting.Report, 0, len(items))
	for _, item := range items {
		if filter.EngagementID != "" && !strings.EqualFold(item.EngagementID, filter.EngagementID) {
			continue
		}
		filtered = append(filtered, item)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].GeneratedAt.After(filtered[j].GeneratedAt) })
	page, size := filter.Page, filter.PageSize
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 25
	}
	if size > 100 {
		size = 100
	}
	start := (page - 1) * size
	if start >= len(filtered) {
		return Page[reporting.Report]{Items: []reporting.Report{}, Page: page, PageSize: size, Total: len(filtered)}
	}
	end := start + size
	if end > len(filtered) {
		end = len(filtered)
	}
	return Page[reporting.Report]{Items: filtered[start:end], Page: page, PageSize: size, Total: len(filtered), HasNext: end < len(filtered)}
}
