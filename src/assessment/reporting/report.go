package reporting

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"ANGEL/src/assessment/checks"
)

type Finding struct {
	ID          string   `json:"id"`
	CheckID     string   `json:"check_id"`
	Title       string   `json:"title"`
	Severity    string   `json:"severity"`
	Confidence  float64  `json:"confidence"`
	Target      string   `json:"target"`
	Summary     string   `json:"summary"`
	EvidenceIDs []string `json:"evidence_ids"`
	Remediation string   `json:"remediation"`
	Status      string   `json:"status"`
}
type Report struct {
	ID            string    `json:"id"`
	EngagementID  string    `json:"engagement_id"`
	GeneratedAt   time.Time `json:"generated_at"`
	EngineVersion string    `json:"engine_version"`
	Findings      []Finding `json:"findings"`
	IntegrityHash string    `json:"integrity_hash"`
}

func New(engagementID string, findings []Finding, at time.Time) (Report, error) {
	if strings.TrimSpace(engagementID) == "" {
		return Report{}, fmt.Errorf("engagement id is required")
	}
	copyFindings := append([]Finding(nil), findings...)
	sort.Slice(copyFindings, func(i, j int) bool { return copyFindings[i].ID < copyFindings[j].ID })
	report := Report{EngagementID: engagementID, GeneratedAt: at.UTC(), EngineVersion: "ANGEL-25", Findings: copyFindings}
	raw, _ := json.Marshal(report)
	hash := sha256.Sum256(raw)
	report.ID = hex.EncodeToString(hash[:])[:24]
	raw, _ = json.Marshal(report)
	hash = sha256.Sum256(raw)
	report.IntegrityHash = hex.EncodeToString(hash[:])
	return report, nil
}
func (r Report) Validate() error {
	if r.ID == "" || r.EngagementID == "" || r.IntegrityHash == "" {
		return fmt.Errorf("report identity and integrity hash are required")
	}
	for _, finding := range r.Findings {
		if finding.ID == "" || finding.CheckID == "" || finding.Title == "" || finding.Target == "" {
			return fmt.Errorf("finding identity is incomplete")
		}
		if finding.Confidence < 0 || finding.Confidence > 1 {
			return fmt.Errorf("finding confidence is invalid")
		}
		if finding.Status == "" {
			return fmt.Errorf("finding status is required")
		}
	}
	if !r.VerifyIntegrity() {
		return fmt.Errorf("report integrity hash mismatch")
	}
	return nil
}
func (r Report) VerifyIntegrity() bool {
	expected := r.IntegrityHash
	r.IntegrityHash = ""
	raw, err := json.Marshal(r)
	if err != nil {
		return false
	}
	hash := sha256.Sum256(raw)
	return expected == hex.EncodeToString(hash[:])
}
func (r Report) JSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(r, "", "  ")
}
func (r Report) Markdown() string {
	var b strings.Builder
	b.WriteString("# ANGEL-25 Security Validation Report\n\n")
	b.WriteString(fmt.Sprintf("Engagement: `%s`\n\nGenerated: `%s`\n\n", r.EngagementID, r.GeneratedAt.Format(time.RFC3339)))
	b.WriteString("## Findings\n\n")
	for _, f := range r.Findings {
		b.WriteString(fmt.Sprintf("### %s — %s\n\n- **Severity:** %s\n- **Confidence:** %.2f\n- **Target:** `%s`\n- **Status:** %s\n\n%s\n\n**Remediation:** %s\n\n", f.ID, f.Title, f.Severity, f.Confidence, f.Target, f.Status, f.Summary, f.Remediation))
	}
	return b.String()
}
func FromCheck(check checks.Result, target, evidenceID string) Finding {
	return Finding{ID: check.CheckID + "-" + target, CheckID: check.CheckID, Title: check.Title, Severity: string(check.Severity), Confidence: check.Confidence, Target: target, Summary: check.Summary, EvidenceIDs: []string{evidenceID}, Remediation: check.Remediation, Status: func() string {
		if check.Passed {
			return "ACCEPTED"
		}
		return "OPEN"
	}()}
}
