package replay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Task struct {
	ID          string  `json:"id"`
	AgentType   string  `json:"agent_type"`
	TargetRef   string  `json:"target_ref"`
	Technique   string  `json:"technique"`
	Mode        string  `json:"mode"`
	RequestedBy string  `json:"requested_by"`
	ApprovalID  *string `json:"approval_id"`
}

type Result struct {
	TaskID       string   `json:"task_id"`
	Status       string   `json:"status"`
	Simulation   bool     `json:"simulation"`
	Authorized   bool     `json:"authorized"`
	Message      string   `json:"message,omitempty"`
	EvidenceRefs []string `json:"evidence_refs"`
}

func Replay(taskData, resultData []byte) (Result, error) {
	var task Task
	if err := decodeStrict(taskData, &task); err != nil {
		return Result{}, fmt.Errorf("decode task: %w", err)
	}
	var result Result
	if err := decodeStrict(resultData, &result); err != nil {
		return Result{}, fmt.Errorf("decode result: %w", err)
	}
	if err := validateTask(task); err != nil {
		return Result{}, err
	}
	if err := validateResult(task, result); err != nil {
		return Result{}, err
	}
	return result, nil
}

func decodeStrict(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("multiple JSON values are not allowed")
		}
		return err
	}
	return nil
}

func validateTask(task Task) error {
	if strings.TrimSpace(task.ID) == "" || strings.TrimSpace(task.AgentType) == "" || strings.TrimSpace(task.Technique) == "" || strings.TrimSpace(task.RequestedBy) == "" {
		return errors.New("task identity, agent type, technique, and requester are required")
	}
	if !isFixtureRef(task.TargetRef) {
		return errors.New("task target must use a fixture reference")
	}
	if task.Mode != "observe" && task.Mode != "simulate" {
		return fmt.Errorf("unsupported task mode %q", task.Mode)
	}
	return nil
}

func validateResult(task Task, result Result) error {
	if result.TaskID != task.ID {
		return errors.New("result task ID does not match task")
	}
	if !result.Simulation {
		return errors.New("replay result must be marked as simulation")
	}
	if !result.Authorized {
		return errors.New("replay result must be authorized")
	}
	switch result.Status {
	case "accepted", "completed", "failed", "rejected":
	default:
		return fmt.Errorf("unsupported result status %q", result.Status)
	}
	for _, ref := range result.EvidenceRefs {
		if !isFixtureRef(ref) {
			return errors.New("evidence references must use fixture references")
		}
	}
	return nil
}

func isFixtureRef(value string) bool {
	return strings.HasPrefix(strings.TrimSpace(value), "fixture://") && len(strings.TrimSpace(value)) > len("fixture://")
}
