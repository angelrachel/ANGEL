package normalization

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strings"
	"time"

	"ANGEL/src/assessment/collectors"
)

type Asset struct {
	ID         string    `json:"id"`
	Canonical  string    `json:"canonical"`
	Kind       string    `json:"kind"`
	Tags       []string  `json:"tags"`
	Confidence float64   `json:"confidence"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
}
type Pipeline struct{ now func() time.Time }

func New() Pipeline { return Pipeline{now: func() time.Time { return time.Now().UTC() }} }
func (p Pipeline) Asset(target, kind string, tags []string, confidence float64) (Asset, error) {
	canonical, err := CanonicalTarget(target)
	if err != nil {
		return Asset{}, err
	}
	if strings.TrimSpace(kind) == "" {
		return Asset{}, fmt.Errorf("asset kind is required")
	}
	if confidence < 0 || confidence > 1 {
		return Asset{}, fmt.Errorf("confidence must be between 0 and 1")
	}
	unique := map[string]bool{}
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag != "" {
			unique[tag] = true
		}
	}
	normalized := make([]string, 0, len(unique))
	for tag := range unique {
		normalized = append(normalized, tag)
	}
	sort.Strings(normalized)
	now := p.now()
	sum := sha256.Sum256([]byte(kind + "|" + canonical))
	return Asset{ID: hex.EncodeToString(sum[:])[:24], Canonical: canonical, Kind: kind, Tags: normalized, Confidence: confidence, FirstSeen: now, LastSeen: now}, nil
}
func CanonicalTarget(target string) (string, error) {
	value := strings.TrimSpace(strings.ToLower(target))
	if value == "" {
		return "", fmt.Errorf("target is required")
	}
	if strings.HasPrefix(value, "fixture://") {
		return value, nil
	}
	if ip := net.ParseIP(value); ip != nil {
		return ip.String(), nil
	}
	if u, err := url.Parse(value); err == nil && u.Hostname() != "" && (u.Scheme == "http" || u.Scheme == "https") {
		return u.Scheme + "://" + strings.ToLower(u.Hostname()), nil
	}
	if strings.Contains(value, ".") && !strings.ContainsAny(value, "/ :") {
		return strings.TrimSuffix(value, "."), nil
	}
	return "", fmt.Errorf("unsupported target format")
}
func MergeObservation(existing Asset, observation collectors.Observation) Asset {
	if existing.Canonical == "" {
		existing.Canonical = observation.Target
	}
	if observation.Confidence > existing.Confidence {
		existing.Confidence = observation.Confidence
	}
	existing.LastSeen = observation.CollectedAt
	existing.Tags = append(existing.Tags, observation.Kind)
	sort.Strings(existing.Tags)
	return existing
}
