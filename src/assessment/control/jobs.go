package control

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"ANGEL/src/assessment/domain"
	"ANGEL/src/assessment/policy"
)

const (
	Created   = "CREATED"
	Approved  = "APPROVED"
	Queued    = "QUEUED"
	Running   = "RUNNING"
	Completed = "COMPLETED"
	Failed    = "FAILED"
	Cancelled = "CANCELLED"
	Expired   = "EXPIRED"
)

type Controller struct {
	mu      sync.RWMutex
	jobs    map[string]domain.AssessmentJob
	policy  *policy.Engine
	private ed25519.PrivateKey
	public  ed25519.PublicKey
}

func NewController(p *policy.Engine) (*Controller, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Controller{jobs: make(map[string]domain.AssessmentJob), policy: p, private: priv, public: pub}, nil
}

func (c *Controller) PublicKey() ed25519.PublicKey {
	return append(ed25519.PublicKey(nil), c.public...)
}

func (c *Controller) Create(engagementID, taskType, target, action string, now time.Time) (domain.AssessmentJob, error) {
	if strings.TrimSpace(taskType) == "" {
		return domain.AssessmentJob{}, fmt.Errorf("task type is required")
	}
	decision := c.policy.Authorize(engagementID, target, action, now)
	if !decision.Allowed {
		return domain.AssessmentJob{}, fmt.Errorf("job denied: %s", decision.Reason)
	}
	idHash := sha256.Sum256([]byte(engagementID + taskType + target + now.UTC().String()))
	id := hex.EncodeToString(idHash[:])[:24]
	job := domain.AssessmentJob{ID: id, EngagementID: engagementID, TaskType: taskType, Target: target, Status: Created, CreatedAt: now.UTC()}
	payload, _ := json.Marshal(job)
	job.Signature = hex.EncodeToString(ed25519.Sign(c.private, payload))
	c.mu.Lock()
	c.jobs[id] = job
	c.mu.Unlock()
	return job, nil
}

func (c *Controller) Transition(id, status string, at time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	job, ok := c.jobs[id]
	if !ok {
		return fmt.Errorf("job not found")
	}
	if !legal(job.Status, status) {
		return fmt.Errorf("invalid transition %s -> %s", job.Status, status)
	}
	job.Status = status
	if status == Running {
		job.StartedAt = at.UTC()
	}
	if status == Completed || status == Failed || status == Cancelled || status == Expired {
		job.FinishedAt = at.UTC()
	}
	c.jobs[id] = job
	return nil
}

func (c *Controller) Get(id string) (domain.AssessmentJob, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	job, ok := c.jobs[id]
	return job, ok
}
func (c *Controller) List() []domain.AssessmentJob {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]domain.AssessmentJob, 0, len(c.jobs))
	for _, job := range c.jobs {
		out = append(out, job)
	}
	return out
}

func legal(from, to string) bool {
	switch from {
	case Created:
		return to == Approved || to == Cancelled || to == Expired
	case Approved:
		return to == Queued || to == Cancelled || to == Expired
	case Queued:
		return to == Running || to == Cancelled || to == Expired
	case Running:
		return to == Completed || to == Failed || to == Cancelled
	default:
		return false
	}
}
