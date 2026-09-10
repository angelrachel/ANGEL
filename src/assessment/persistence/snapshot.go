package persistence

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"ANGEL/src/assessment/domain"
	"ANGEL/src/assessment/evidence"
	"ANGEL/src/assessment/governance"
	"ANGEL/src/assessment/reporting"
)

type Snapshot struct {
	Version     int                    `json:"version"`
	ExportedAt  time.Time              `json:"exported_at"`
	Engagements []domain.Engagement    `json:"engagements"`
	Jobs        []domain.AssessmentJob `json:"jobs"`
	Findings    []reporting.Finding    `json:"findings"`
	Evidence    []evidence.Bundle      `json:"evidence"`
	Audit       []governance.Event     `json:"audit"`
	Digest      string                 `json:"digest"`
	Signature   string                 `json:"signature"`
}
type Signer struct {
	Private ed25519.PrivateKey
	Public  ed25519.PublicKey
}

func NewSigner() (Signer, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	return Signer{Private: priv, Public: pub}, err
}
func Create(snapshot Snapshot, signer Signer) (Snapshot, error) {
	if snapshot.Version < 1 {
		snapshot.Version = 1
	}
	if snapshot.ExportedAt.IsZero() {
		snapshot.ExportedAt = time.Now().UTC()
	}
	snapshot.Signature = ""
	snapshot.Digest = ""
	raw, err := json.Marshal(snapshot)
	if err != nil {
		return Snapshot{}, err
	}
	sum := sha256.Sum256(raw)
	snapshot.Digest = hex.EncodeToString(sum[:])
	raw, err = json.Marshal(snapshot)
	if err != nil {
		return Snapshot{}, err
	}
	snapshot.Signature = hex.EncodeToString(ed25519.Sign(signer.Private, raw))
	return snapshot, nil
}
func Verify(snapshot Snapshot, public ed25519.PublicKey) error {
	if snapshot.Version < 1 || strings.TrimSpace(snapshot.Digest) == "" || strings.TrimSpace(snapshot.Signature) == "" {
		return fmt.Errorf("snapshot identity is incomplete")
	}
	signature := snapshot.Signature
	digest := snapshot.Digest
	snapshot.Signature = ""
	snapshot.Digest = ""
	raw, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != digest {
		return fmt.Errorf("snapshot digest mismatch")
	}
	snapshot.Digest = digest
	raw, err = json.Marshal(snapshot)
	if err != nil {
		return err
	}
	decoded, err := hex.DecodeString(signature)
	if err != nil || !ed25519.Verify(public, raw, decoded) {
		return fmt.Errorf("snapshot signature mismatch")
	}
	return nil
}
func Encode(snapshot Snapshot) ([]byte, error) {
	if strings.TrimSpace(snapshot.Digest) == "" || strings.TrimSpace(snapshot.Signature) == "" {
		return nil, fmt.Errorf("unsigned snapshot")
	}
	return json.MarshalIndent(snapshot, "", "  ")
}
func Decode(data []byte) (Snapshot, error) {
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return snapshot, err
	}
	return snapshot, nil
}
