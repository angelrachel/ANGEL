package evidence

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Bundle struct {
	ID           string    `json:"id"`
	EngagementID string    `json:"engagement_id"`
	JobID        string    `json:"job_id"`
	CapturedAt   time.Time `json:"captured_at"`
	ContentType  string    `json:"content_type"`
	SHA256       string    `json:"sha256"`
	ParentSHA256 string    `json:"parent_sha256"`
	Redacted     bool      `json:"redacted"`
	Payload      []byte    `json:"payload"`
	Signature    string    `json:"signature"`
}
type Signer struct {
	Private ed25519.PrivateKey
	Public  ed25519.PublicKey
}

func NewSigner() (Signer, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	return Signer{Private: priv, Public: pub}, err
}
func CreateBundle(engagementID, jobID, contentType string, payload []byte, parent string, redacted bool, at time.Time, signer Signer) (Bundle, error) {
	if strings.TrimSpace(engagementID) == "" || strings.TrimSpace(jobID) == "" || len(payload) == 0 {
		return Bundle{}, fmt.Errorf("engagement, job, and payload are required")
	}
	hash := sha256.Sum256(payload)
	id := hex.EncodeToString(hash[:])[:24]
	bundle := Bundle{ID: id, EngagementID: engagementID, JobID: jobID, CapturedAt: at.UTC(), ContentType: contentType, SHA256: hex.EncodeToString(hash[:]), ParentSHA256: parent, Redacted: redacted, Payload: append([]byte(nil), payload...)}
	raw, _ := json.Marshal(bundle)
	bundle.Signature = hex.EncodeToString(ed25519.Sign(signer.Private, raw))
	return bundle, nil
}
func VerifyBundle(bundle Bundle, public ed25519.PublicKey) bool {
	hash := sha256.Sum256(bundle.Payload)
	if hex.EncodeToString(hash[:]) != bundle.SHA256 || bundle.ID != hex.EncodeToString(hash[:])[:24] {
		return false
	}
	signature, err := hex.DecodeString(bundle.Signature)
	if err != nil {
		return false
	}
	copyBundle := bundle
	copyBundle.Signature = ""
	raw, _ := json.Marshal(copyBundle)
	return ed25519.Verify(public, raw, signature)
}
