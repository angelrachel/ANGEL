package storage

import (
	"sort"
	"strings"

	"ANGEL/src/assessment/domain"
)

type ObservationQuery struct {
	JobID    string
	CheckID  string
	Status   string
	Page     int
	PageSize int
}

type ObservationPage struct {
	Items    []domain.Observation `json:"items"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
	Total    int                  `json:"total"`
	HasNext  bool                 `json:"has_next"`
}

func QueryObservations(items []domain.Observation, query ObservationQuery) ObservationPage {
	filtered := make([]domain.Observation, 0, len(items))
	for _, item := range items {
		if query.JobID != "" && item.JobID != query.JobID {
			continue
		}
		if query.CheckID != "" && item.CheckID != query.CheckID {
			continue
		}
		if query.Status != "" && !strings.EqualFold(item.Status, query.Status) {
			continue
		}
		filtered = append(filtered, item)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].CreatedAt.After(filtered[j].CreatedAt) })
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
		return ObservationPage{Items: []domain.Observation{}, Page: page, PageSize: size, Total: total}
	}
	end := start + size
	if end > total {
		end = total
	}
	return ObservationPage{Items: filtered[start:end], Page: page, PageSize: size, Total: total, HasNext: end < total}
}
