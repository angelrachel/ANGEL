package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"ANGEL/src/assessment/domain"
	"ANGEL/src/assessment/evidence"
	"ANGEL/src/assessment/reporting"
)

type PersistentState struct {
	Version           int                    `json:"version"`
	UpdatedAt         time.Time              `json:"updated_at"`
	ControllerPrivate string                 `json:"controller_private,omitempty"`
	ControllerPublic  string                 `json:"controller_public,omitempty"`
	Jobs              []domain.AssessmentJob `json:"jobs"`
	Observations      []domain.Observation   `json:"observations"`
	Evidence          []evidence.Bundle      `json:"evidence"`
	Findings          []reporting.Finding    `json:"findings"`
	Reports           []reporting.Report     `json:"reports"`
	Digest            string                 `json:"digest"`
}

type FileStore struct {
	mu   sync.Mutex
	path string
}

func OpenFileStore(path string) (*FileStore, PersistentState, error) {
	if path == "" {
		return nil, PersistentState{}, fmt.Errorf("state path is required")
	}
	store := &FileStore{path: path}
	state, err := store.Load()
	if os.IsNotExist(err) {
		return store, PersistentState{Version: 1}, nil
	}
	if err != nil {
		return nil, PersistentState{}, err
	}
	return store, state, nil
}

func (s *FileStore) Load() (PersistentState, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return PersistentState{}, err
	}
	var envelope PersistentState
	if err := json.Unmarshal(data, &envelope); err != nil {
		return envelope, fmt.Errorf("decode durable state: %w", err)
	}
	if envelope.Version < 1 || envelope.Digest == "" {
		return envelope, fmt.Errorf("durable state identity is incomplete")
	}
	expected := envelope.Digest
	envelope.Digest = ""
	raw, err := json.Marshal(envelope)
	if err != nil {
		return envelope, err
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != expected {
		return envelope, fmt.Errorf("durable state digest mismatch")
	}
	envelope.Digest = expected
	return envelope, nil
}

func (s *FileStore) Save(state PersistentState) error {
	if s == nil || s.path == "" {
		return fmt.Errorf("state store is not configured")
	}
	if state.Version < 1 {
		state.Version = 1
	}
	state.UpdatedAt = time.Now().UTC()
	state.Digest = ""
	raw, err := json.Marshal(state)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(raw)
	state.Digest = hex.EncodeToString(sum[:])
	raw, err = json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(s.path), ".angel-state-*")
	if err != nil {
		return err
	}
	tempName := temp.Name()
	defer os.Remove(tempName)
	if err := temp.Chmod(0600); err != nil {
		temp.Close()
		return err
	}
	if _, err := temp.Write(raw); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Sync(); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(tempName, s.path)
}

func (s *FileStore) Path() string {
	if s == nil {
		return ""
	}
	return s.path
}
