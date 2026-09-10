package assets

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"ANGEL/src/assessment/domain"
)

var criticalities = map[string]bool{"critical": true, "high": true, "medium": true, "low": true, "unknown": true}

// Normalize converts an authorized asset reference into a stable representation.
// It never performs network access or DNS resolution.
func Normalize(kind, value string) (string, error) {
	kind = strings.ToLower(strings.TrimSpace(kind))
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("asset value is required")
	}
	switch kind {
	case "hostname":
		value = strings.TrimSuffix(strings.ToLower(value), ".")
		if strings.ContainsAny(value, "/:@") || net.ParseIP(value) != nil {
			return "", fmt.Errorf("invalid hostname asset")
		}
		return value, nil
	case "url":
		parsed, err := url.Parse(value)
		if err != nil || parsed.Hostname() == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return "", fmt.Errorf("invalid URL asset")
		}
		parsed.Scheme = strings.ToLower(parsed.Scheme)
		parsed.Host = strings.ToLower(parsed.Host)
		parsed.Fragment = ""
		return strings.TrimSuffix(parsed.String(), "/"), nil
	case "cidr":
		ip, block, err := net.ParseCIDR(value)
		if err != nil || ip == nil {
			return "", fmt.Errorf("invalid CIDR asset")
		}
		return block.String(), nil
	case "fixture":
		if !strings.HasPrefix(value, "fixture://") {
			return "", fmt.Errorf("fixture asset must use fixture://")
		}
		return value, nil
	default:
		return "", fmt.Errorf("unsupported asset kind: %s", kind)
	}
}

func New(engagementID, kind, value, owner, criticality string, confidence float64, at time.Time) (domain.Asset, error) {
	if strings.TrimSpace(engagementID) == "" {
		return domain.Asset{}, fmt.Errorf("engagement ID is required")
	}
	canonical, err := Normalize(kind, value)
	if err != nil {
		return domain.Asset{}, err
	}
	criticality = strings.ToLower(strings.TrimSpace(criticality))
	if !criticalities[criticality] {
		return domain.Asset{}, fmt.Errorf("unsupported criticality: %s", criticality)
	}
	if confidence < 0 || confidence > 1 {
		return domain.Asset{}, fmt.Errorf("confidence must be between 0 and 1")
	}
	digest := sha256.Sum256([]byte(engagementID + "|" + kind + "|" + canonical))
	return domain.Asset{ID: hex.EncodeToString(digest[:])[:24], EngagementID: engagementID, Canonical: canonical, Kind: strings.ToLower(strings.TrimSpace(kind)), Owner: strings.TrimSpace(owner), Criticality: criticality, Confidence: confidence, LastSeen: at.UTC()}, nil
}
