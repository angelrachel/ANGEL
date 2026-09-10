package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

// Manifest is a deterministic summary of a verified evidence chain.
type Manifest struct {
	RecordCount int    `json:"record_count"`
	FirstID     string `json:"first_id,omitempty"`
	LastID      string `json:"last_id,omitempty"`
	ChainHash   string `json:"chain_hash"`
}

// BuildManifest validates records before producing a shareable chain summary.
func BuildManifest(records []Record) (Manifest, error) {
	if !VerifyChain(records) {
		return Manifest{}, errors.New("cannot build manifest from invalid evidence chain")
	}
	manifest := Manifest{RecordCount: len(records), ChainHash: chainHash(records)}
	if len(records) > 0 {
		manifest.FirstID = records[0].ID
		manifest.LastID = records[len(records)-1].ID
	}
	return manifest, nil
}

func chainHash(records []Record) string {
	encoded, _ := json.Marshal(records)
	hash := sha256.Sum256(encoded)
	return hex.EncodeToString(hash[:])
}

func (m Manifest) Validate(records []Record) error {
	if !VerifyChain(records) {
		return errors.New("evidence chain is invalid")
	}
	expected, err := BuildManifest(records)
	if err != nil {
		return err
	}
	if m != expected {
		return fmt.Errorf("manifest mismatch: got %#v, expected %#v", m, expected)
	}
	return nil
}
