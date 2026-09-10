package storage

import (
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
)

func TestJobRepositoryUsesOptimisticVersioning(t *testing.T) {
	repo := NewJobRepository()
	job := domain.AssessmentJob{ID: "job-1", CreatedAt: time.Now()}
	if err := repo.Put(job, 0); err != nil {
		t.Fatal(err)
	}
	if err := repo.Put(job, 0); err == nil {
		t.Fatal("stale version accepted")
	}
	_, version, ok := repo.Get(job.ID)
	if !ok || version != 1 {
		t.Fatalf("version=%d ok=%v", version, ok)
	}
}
func TestObservationRepositoryIndexesByJob(t *testing.T) {
	repo := NewObservationRepository()
	if err := repo.Put(domain.Observation{ID: "o-1", JobID: "job-1", CheckID: "check", CreatedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if len(repo.ByJob("job-1")) != 1 {
		t.Fatal("observation missing")
	}
}
func TestRetentionPolicy(t *testing.T) {
	policy := RetentionPolicy{EvidenceDays: 7}
	if !policy.Expired(time.Now().Add(-8*24*time.Hour), time.Now()) {
		t.Fatal("expired evidence not detected")
	}
}
