package storage

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"ANGEL/src/assessment/domain"
)

type JobRepository struct {
	mu       sync.RWMutex
	items    map[string]domain.AssessmentJob
	versions map[string]int64
}

func NewJobRepository() *JobRepository {
	return &JobRepository{items: map[string]domain.AssessmentJob{}, versions: map[string]int64{}}
}
func (r *JobRepository) Put(job domain.AssessmentJob, expectedVersion int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	version := r.versions[job.ID]
	if version != expectedVersion {
		return fmt.Errorf("version conflict: expected %d got %d", expectedVersion, version)
	}
	r.items[job.ID] = job
	r.versions[job.ID] = version + 1
	return nil
}
func (r *JobRepository) Get(id string) (domain.AssessmentJob, int64, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	job, ok := r.items[id]
	return job, r.versions[id], ok
}
func (r *JobRepository) List(limit int) []domain.AssessmentJob {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.AssessmentJob, 0, len(r.items))
	for _, item := range r.items {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

type ObservationRepository struct {
	mu    sync.RWMutex
	items map[string]domain.Observation
}

func NewObservationRepository() *ObservationRepository {
	return &ObservationRepository{items: map[string]domain.Observation{}}
}
func (r *ObservationRepository) Put(item domain.Observation) error {
	if item.ID == "" || item.JobID == "" || item.CheckID == "" {
		return fmt.Errorf("observation identity is required")
	}
	r.mu.Lock()
	r.items[item.ID] = item
	r.mu.Unlock()
	return nil
}
func (r *ObservationRepository) ByJob(jobID string) []domain.Observation {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []domain.Observation{}
	for _, item := range r.items {
		if item.JobID == jobID {
			out = append(out, item)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}

type RetentionPolicy struct {
	EvidenceDays int
	ReportDays   int
	AuditDays    int
}

func (p RetentionPolicy) Expired(captured, now time.Time) bool {
	return p.EvidenceDays > 0 && now.After(captured.Add(time.Duration(p.EvidenceDays)*24*time.Hour))
}
