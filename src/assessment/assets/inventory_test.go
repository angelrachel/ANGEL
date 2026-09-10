package assets

import (
	"testing"
	"time"
)

func TestNormalizeAssetReferences(t *testing.T) {
	cases := []struct{ kind, value, want string }{
		{"hostname", "Example.COM.", "example.com"},
		{"url", "HTTPS://Example.COM/path/#fragment", "https://example.com/path"},
		{"cidr", "192.0.2.1/24", "192.0.2.0/24"},
		{"fixture", "fixture://lab/web", "fixture://lab/web"},
	}
	for _, tc := range cases {
		got, err := Normalize(tc.kind, tc.value)
		if err != nil || got != tc.want {
			t.Fatalf("Normalize(%q,%q)=%q,%v want %q", tc.kind, tc.value, got, err, tc.want)
		}
	}
}

func TestNewAssetRejectsInvalidInput(t *testing.T) {
	if _, err := New("eng", "hostname", "", "owner", "high", .9, time.Now()); err == nil {
		t.Fatal("empty asset accepted")
	}
	if _, err := New("eng", "hostname", "example.com", "owner", "urgent", .9, time.Now()); err == nil {
		t.Fatal("invalid criticality accepted")
	}
	if _, err := New("eng", "hostname", "example.com", "owner", "high", 2, time.Now()); err == nil {
		t.Fatal("invalid confidence accepted")
	}
}

func TestNewAssetIdentityIsStable(t *testing.T) {
	at := time.Unix(100, 0)
	first, err := New("eng", "hostname", "Example.COM.", "owner", "high", .9, at)
	if err != nil {
		t.Fatal(err)
	}
	second, err := New("eng", "hostname", "example.com", "different-owner", "high", .9, at)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != second.ID {
		t.Fatalf("asset identity changed: %q %q", first.ID, second.ID)
	}
}
