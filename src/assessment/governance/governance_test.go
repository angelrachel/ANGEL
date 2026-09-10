package governance

import (
	"testing"
	"time"
)

func TestAuditLogVerifiesAndDetectsTampering(t *testing.T) {
	log := NewAuditLog()
	log.Append("operator", "JOB_CREATED", "job", "job-1", time.Now())
	log.Append("worker", "EVIDENCE_ADDED", "evidence", "ev-1", time.Now())
	events := log.List()
	if !Verify(events) {
		t.Fatal("valid audit chain rejected")
	}
	events[1].Action = "JOB_DELETED"
	if Verify(events) {
		t.Fatal("tampered audit chain accepted")
	}
}
func TestProductionConfigRequiresStrongOperatorKey(t *testing.T) {
	config := Config{Environment: "production", DatabaseURL: "postgres://db", EvidenceBucket: "s3://evidence", OperatorKey: "short", RetentionDays: 30}
	if err := config.Validate(); err == nil {
		t.Fatal("weak production key accepted")
	}
	config.OperatorKey = "strong-production-operator-key"
	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
}
