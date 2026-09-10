package main

import (
	"encoding/hex"

	"ANGEL/src/assessment/domain"
	"ANGEL/src/assessment/evidence"
	"ANGEL/src/assessment/reporting"
	"ANGEL/src/assessment/storage"
)

func (e *Engine) persistState() {
	if e == nil || e.durableState == nil {
		return
	}
	e.mu.RLock()
	evidence := make([]evidence.Bundle, 0, len(e.evidence))
	for _, item := range e.evidence {
		evidence = append(evidence, item)
	}
	observations := make([]domain.Observation, 0, len(e.observations))
	for _, item := range e.observations {
		observations = append(observations, item)
	}
	findings := make([]reporting.Finding, 0, len(e.findings))
	for _, item := range e.findings {
		findings = append(findings, item)
	}
	reports := make([]reporting.Report, 0, len(e.reports))
	for _, item := range e.reports {
		reports = append(reports, item)
	}
	e.mu.RUnlock()
	private, public := e.jobController.KeyMaterial()
	state := storage.PersistentState{Version: 1, Jobs: e.jobController.List(), Policy: e.policyEngine.Snapshot(), Observations: observations, Evidence: evidence, Findings: findings, Reports: reports, ControllerPrivate: hex.EncodeToString(private), ControllerPublic: hex.EncodeToString(public)}
	_ = e.durableState.Save(state)
}
