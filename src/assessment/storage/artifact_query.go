package storage

import (
	"sort"
	"strings"

	"ANGEL/src/assessment/evidence"
	"ANGEL/src/assessment/reporting"
)

type EvidenceQuery struct {
	EngagementID string
	JobID        string
	ContentType  string
	Page         int
	PageSize     int
}
type EvidencePage struct {
	Items    []evidence.Bundle `json:"items"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Total    int               `json:"total"`
	HasNext  bool              `json:"has_next"`
}

func QueryEvidence(items []evidence.Bundle, query EvidenceQuery) EvidencePage {
	filtered := make([]evidence.Bundle, 0, len(items))
	for _, item := range items {
		if query.EngagementID != "" && item.EngagementID != query.EngagementID {
			continue
		}
		if query.JobID != "" && item.JobID != query.JobID {
			continue
		}
		if query.ContentType != "" && !strings.EqualFold(item.ContentType, query.ContentType) {
			continue
		}
		safe := item
		safe.Payload = nil
		filtered = append(filtered, safe)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].CapturedAt.After(filtered[j].CapturedAt) })
	page := query.Page
	if page < 1 {
		page = 1
	}
	size := query.PageSize
	if size < 1 || size > 100 {
		size = 25
	}
	total := len(filtered)
	start := (page - 1) * size
	if start >= total {
		return EvidencePage{Items: []evidence.Bundle{}, Page: page, PageSize: size, Total: total}
	}
	end := start + size
	if end > total {
		end = total
	}
	return EvidencePage{Items: filtered[start:end], Page: page, PageSize: size, Total: total, HasNext: end < total}
}

type ReportQuery struct {
	EngagementID string
	Page         int
	PageSize     int
}
type ReportPage struct {
	Items    []reporting.Report `json:"items"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
	Total    int                `json:"total"`
	HasNext  bool               `json:"has_next"`
}

func QueryReports(items []reporting.Report, query ReportQuery) ReportPage {
	filtered := make([]reporting.Report, 0, len(items))
	for _, item := range items {
		if query.EngagementID != "" && item.EngagementID != query.EngagementID {
			continue
		}
		filtered = append(filtered, item)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].GeneratedAt.After(filtered[j].GeneratedAt) })
	page := query.Page
	if page < 1 {
		page = 1
	}
	size := query.PageSize
	if size < 1 || size > 100 {
		size = 25
	}
	total := len(filtered)
	start := (page - 1) * size
	if start >= total {
		return ReportPage{Items: []reporting.Report{}, Page: page, PageSize: size, Total: total}
	}
	end := start + size
	if end > total {
		end = total
	}
	return ReportPage{Items: filtered[start:end], Page: page, PageSize: size, Total: total, HasNext: end < total}
}
