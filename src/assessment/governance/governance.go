package governance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Event struct {
	ID         string    `json:"id"`
	Actor      string    `json:"actor"`
	Action     string    `json:"action"`
	Object     string    `json:"object"`
	ObjectID   string    `json:"object_id"`
	At         time.Time `json:"at"`
	ParentHash string    `json:"parent_hash"`
	Hash       string    `json:"hash"`
}
type AuditLog struct {
	mu     sync.RWMutex
	events []Event
}

func NewAuditLog() *AuditLog { return &AuditLog{events: []Event{}} }
func (l *AuditLog) Append(actor, action, object, objectID string, at time.Time) Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	parent := "0"
	if len(l.events) > 0 {
		parent = l.events[len(l.events)-1].Hash
	}
	payload := fmt.Sprintf("%s|%s|%s|%s|%s|%s", actor, action, object, objectID, at.UTC().Format(time.RFC3339Nano), parent)
	sum := sha256.Sum256([]byte(payload))
	hash := hex.EncodeToString(sum[:])
	event := Event{ID: hash[:24], Actor: actor, Action: action, Object: object, ObjectID: objectID, At: at.UTC(), ParentHash: parent, Hash: hash}
	l.events = append(l.events, event)
	return event
}
func (l *AuditLog) List() []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return append([]Event(nil), l.events...)
}
func Verify(events []Event) bool {
	parent := "0"
	for _, event := range events {
		payload := fmt.Sprintf("%s|%s|%s|%s|%s|%s", event.Actor, event.Action, event.Object, event.ObjectID, event.At.UTC().Format(time.RFC3339Nano), parent)
		sum := sha256.Sum256([]byte(payload))
		hash := hex.EncodeToString(sum[:])
		if event.ParentHash != parent || event.Hash != hash || event.ID != hash[:24] {
			return false
		}
		parent = hash
	}
	return true
}

type Config struct {
	Environment    string
	DatabaseURL    string
	EvidenceBucket string
	OperatorKey    string
	RetentionDays  int
	ReadOnly       bool
}

func (c Config) Validate() error {
	if c.Environment != "development" && c.Environment != "staging" && c.Environment != "production" {
		return fmt.Errorf("unsupported environment")
	}
	if strings.TrimSpace(c.DatabaseURL) == "" {
		return fmt.Errorf("database URL is required")
	}
	if strings.TrimSpace(c.EvidenceBucket) == "" {
		return fmt.Errorf("evidence bucket is required")
	}
	if c.RetentionDays < 1 {
		return fmt.Errorf("retention days must be positive")
	}
	if c.Environment == "production" && (len(c.OperatorKey) < 16 || strings.Contains(c.OperatorKey, "replace-with")) {
		return fmt.Errorf("production operator key is invalid")
	}
	return nil
}
