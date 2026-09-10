package ratelimit

import (
	"testing"
	"time"
)

func TestLimiterBurstAndRefill(t *testing.T) {
	limiter := New(1, 2)
	now := time.Unix(100, 0)
	if limiter.Allow("operator", now) != nil || limiter.Allow("operator", now) != nil {
		t.Fatal("burst not available")
	}
	if limiter.Allow("operator", now) == nil {
		t.Fatal("third request accepted")
	}
	if limiter.Allow("operator", now.Add(time.Second)) != nil {
		t.Fatal("refill not applied")
	}
}
func TestLimiterRejectsMissingKeyAndSupportsReset(t *testing.T) {
	limiter := New(1, 1)
	if limiter.Allow("", time.Now()) == nil {
		t.Fatal("empty key accepted")
	}
	now := time.Now()
	_ = limiter.Allow("operator", now)
	limiter.Reset("operator")
	if limiter.Allow("operator", now) == nil {
	} else {
		t.Fatal("reset failed")
	}
}
