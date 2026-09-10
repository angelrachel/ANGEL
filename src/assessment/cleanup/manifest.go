package cleanup

import (
	"errors"
	"fmt"
	"strings"
)

type Manifest struct {
	ID      string  `json:"manifest_id"`
	Mode    string  `json:"mode"`
	Entries []Entry `json:"entries"`
}

type Entry struct {
	PathRef string `json:"path_ref"`
	Status  string `json:"status"`
}

func NewSimulationManifest(id string, paths []string) (Manifest, error) {
	manifest := Manifest{ID: strings.TrimSpace(id), Mode: "simulate"}
	for _, path := range paths {
		manifest.Entries = append(manifest.Entries, Entry{PathRef: strings.TrimSpace(path), Status: "planned"})
	}
	if err := manifest.Validate(); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

func (m Manifest) Validate() error {
	if m.ID == "" {
		return errors.New("cleanup manifest ID is required")
	}
	if m.Mode != "simulate" && m.Mode != "verify" {
		return fmt.Errorf("unsupported cleanup manifest mode %q", m.Mode)
	}
	for i, entry := range m.Entries {
		if !strings.HasPrefix(entry.PathRef, "fixture://") || len(entry.PathRef) == len("fixture://") {
			return fmt.Errorf("entry %d must use a fixture reference", i)
		}
		switch entry.Status {
		case "planned", "verified", "skipped", "failed":
		default:
			return fmt.Errorf("entry %d has unsupported status %q", i, entry.Status)
		}
	}
	return nil
}
