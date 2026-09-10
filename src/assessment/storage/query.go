package storage

import (
	"sort"
	"strings"

	"ANGEL/src/assessment/domain"
)

type JobQuery struct {
	EngagementID string
	Status       string
	Target       string
	Page         int
	PageSize     int
}

type JobPage struct {
	Items    []domain.AssessmentJob `json:"items"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Total    int                    `json:"total"`
	HasNext  bool                   `json:"has_next"`
}

func QueryJobs(jobs []domain.AssessmentJob, query JobQuery) JobPage {
	filtered := make([]domain.AssessmentJob, 0, len(jobs))
	for _, job := range jobs {
		if query.EngagementID != "" && job.EngagementID != query.EngagementID {
			continue
		}
		if query.Status != "" && !strings.EqualFold(job.Status, query.Status) {
			continue
		}
		if query.Target != "" && job.Target != query.Target {
			continue
		}
		filtered = append(filtered, job)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].CreatedAt.After(filtered[j].CreatedAt) })
	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 25
	}
	total := len(filtered)
	start := (page - 1) * pageSize
	if start >= total {
		return JobPage{Items: []domain.AssessmentJob{}, Page: page, PageSize: pageSize, Total: total, HasNext: false}
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return JobPage{Items: filtered[start:end], Page: page, PageSize: pageSize, Total: total, HasNext: end < total}
}
