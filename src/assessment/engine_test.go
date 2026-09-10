package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEngineAuthRequiresBearerAndUsesConfiguredKey(t *testing.T) {
	engine := &Engine{apiKey: "secret"}
	for _, tc := range []struct {
		name   string
		header string
		want   bool
	}{
		{"missing", "", false},
		{"wrong scheme", "Basic secret", false},
		{"wrong key", "Bearer other", false},
		{"valid", "Bearer secret", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Authorization", tc.header)
			if got := engine.auth(r) == nil; got != tc.want {
				t.Fatalf("auth success=%v, want %v", got, tc.want)
			}
		})
	}
}

func TestClientIPHandlesIPv4AndIPv6(t *testing.T) {
	if got := clientIP("192.0.2.10:1234"); got != "192.0.2.10" {
		t.Fatalf("got %q", got)
	}
	if got := clientIP("[2001:db8::1]:443"); got != "2001:db8::1" {
		t.Fatalf("got %q", got)
	}
}

func TestDecodeJSONRejectsTrailingValues(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"id":"a"} {"id":"b"}`))
	w := httptest.NewRecorder()
	var value Agent
	if err := decodeJSON(w, r, &value); err == nil {
		t.Fatal("expected trailing JSON to be rejected")
	}
}

func TestDecodeJSONAcceptsSingleValue(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"id":"a"}`))
	w := httptest.NewRecorder()
	var value Agent
	if err := decodeJSON(w, r, &value); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if value.ID != "a" {
		t.Fatalf("got ID %q", value.ID)
	}
}
