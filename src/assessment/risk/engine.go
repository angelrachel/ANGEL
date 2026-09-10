package risk

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"ANGEL/src/assessment/domain"
)

type Input struct {
	Exploitability         int
	Reachability           int
	Confidentiality        int
	Integrity              int
	Availability           int
	BusinessCriticality    int
	BlastRadius            int
	Confidence             float64
	AuthenticationRequired bool
	DetectionCoverage      int
}
type Result struct {
	Severity   string   `json:"severity"`
	Score      float64  `json:"score"`
	Confidence float64  `json:"confidence"`
	Rationale  []string `json:"rationale"`
}

func Score(in Input) (Result, error) {
	if in.Confidence < 0 || in.Confidence > 1 {
		return Result{}, fmt.Errorf("confidence must be between 0 and 1")
	}
	for _, value := range []int{in.Exploitability, in.Reachability, in.Confidentiality, in.Integrity, in.Availability, in.BusinessCriticality, in.BlastRadius, in.DetectionCoverage} {
		if value < 0 || value > 10 {
			return Result{}, fmt.Errorf("risk factors must be between 0 and 10")
		}
	}
	base := float64(in.Exploitability+in.Reachability+in.Confidentiality+in.Integrity+in.Availability+in.BusinessCriticality+in.BlastRadius) / 7
	if in.AuthenticationRequired {
		base -= 0.8
	}
	base += (10 - float64(in.DetectionCoverage)) * 0.15
	base *= (0.5 + in.Confidence/2)
	severity := "LOW"
	if base >= 8.5 {
		severity = "P1"
	}
	if base >= 9.4 && in.Confidence >= .9 && in.BusinessCriticality >= 9 && in.Reachability >= 8 {
		severity = "P0"
	}
	if severity == "LOW" && base >= 4 {
		severity = "MEDIUM"
	}
	if severity == "MEDIUM" && base >= 6.5 {
		severity = "HIGH"
	}
	return Result{Severity: severity, Score: base, Confidence: in.Confidence, Rationale: []string{fmt.Sprintf("weighted risk score %.2f", base), fmt.Sprintf("business criticality %d/10", in.BusinessCriticality), fmt.Sprintf("evidence confidence %.2f", in.Confidence)}}, nil
}

func Finding(title, impact, remediation string, assets, evidence []string, in Input) (domain.Finding, Result, error) {
	result, err := Score(in)
	if err != nil {
		return domain.Finding{}, Result{}, err
	}
	title = strings.TrimSpace(title)
	impact = strings.TrimSpace(impact)
	remediation = strings.TrimSpace(remediation)
	if title == "" {
		return domain.Finding{}, Result{}, fmt.Errorf("finding title is required")
	}
	if impact == "" || remediation == "" {
		return domain.Finding{}, Result{}, fmt.Errorf("finding impact and remediation are required")
	}
	if len(assets) == 0 || len(evidence) == 0 {
		return domain.Finding{}, Result{}, fmt.Errorf("finding requires affected assets and evidence")
	}
	if (result.Severity == "P0" || result.Severity == "P1") && result.Confidence < .85 {
		return domain.Finding{}, Result{}, fmt.Errorf("high severity finding requires confidence of at least 0.85")
	}
	identity := sha256.Sum256([]byte(title + "|" + strings.Join(assets, ",") + "|" + strings.Join(evidence, ",")))
	finding := domain.Finding{ID: hex.EncodeToString(identity[:])[:24], Title: title, Severity: result.Severity, Confidence: result.Confidence, AffectedAssets: append([]string(nil), assets...), Impact: impact, EvidenceIDs: append([]string(nil), evidence...), Remediation: remediation, Status: "OPEN"}
	return finding, result, nil
}
