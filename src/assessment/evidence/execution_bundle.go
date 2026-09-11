package evidence

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"ANGEL/src/assessment/execution"
)

// BundleExecution converts a completed safe adapter result into the immutable
// evidence format used by reports and chain-of-custody verification.
func BundleExecution(output execution.Output, engagementID, parent string, signer Signer, at time.Time) (Bundle, error) {
	if strings.TrimSpace(output.JobID) == "" || strings.TrimSpace(output.CheckID) == "" {
		return Bundle{}, fmt.Errorf("execution job and check IDs are required")
	}
	if output.Status != "COMPLETED" {
		return Bundle{}, fmt.Errorf("only completed executions can become evidence")
	}
	payload, err := json.Marshal(struct {
		JobID       string    `json:"job_id"`
		CheckID     string    `json:"check_id"`
		Status      string    `json:"status"`
		Result      any       `json:"result"`
		Observation any       `json:"observation"`
		Evidence    []byte    `json:"evidence"`
		CollectedAt time.Time `json:"collected_at"`
	}{output.JobID, output.CheckID, output.Status, output.Result, output.Observation, output.Evidence, output.CollectedAt.UTC()})
	if err != nil {
		return Bundle{}, fmt.Errorf("marshal execution evidence: %w", err)
	}
	return CreateBundle(engagementID, output.JobID, "application/vnd.angel.execution+json", payload, parent, true, at, signer)
}
