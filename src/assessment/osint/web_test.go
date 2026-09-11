package osint

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebFingerprintDetectsWordPress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html>wp-content/plugins/test</html>"))
	}))
	defer server.Close()

	w := NewWebFingerprint()
	result := w.Fingerprint(server.URL)
	if result != "WordPress" {
		t.Fatalf("expected WordPress, got %q", result)
	}
}

func TestWebFingerprintDetectsJoomla(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html>Joomla site</html>"))
	}))
	defer server.Close()

	w := NewWebFingerprint()
	result := w.Fingerprint(server.URL)
	if result != "Joomla" {
		t.Fatalf("expected Joomla, got %q", result)
	}
}

func TestWebFingerprintDetectsDrupal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html>Drupal installation</html>"))
	}))
	defer server.Close()

	w := NewWebFingerprint()
	result := w.Fingerprint(server.URL)
	if result != "Drupal" {
		t.Fatalf("expected Drupal, got %q", result)
	}
}

func TestWebFingerprintDetectsLaravel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html>Laravel app</html>"))
	}))
	defer server.Close()

	w := NewWebFingerprint()
	result := w.Fingerprint(server.URL)
	if result != "Laravel" {
		t.Fatalf("expected Laravel, got %q", result)
	}
}

func TestWebFingerprintReturnsUnknownForUnrecognized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html>Plain HTML site</html>"))
	}))
	defer server.Close()

	w := NewWebFingerprint()
	result := w.Fingerprint(server.URL)
	if result != "Unknown" {
		t.Fatalf("expected Unknown, got %q", result)
	}
}

func TestWebFingerprintReturnsEmptyOnConnectionError(t *testing.T) {
	w := NewWebFingerprint()
	result := w.Fingerprint("http://192.0.2.1:9999/nonexistent")
	if result != "" {
		t.Fatalf("expected empty string on connection error, got %q", result)
	}
}

func TestWebFingerprintGetServerHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "nginx/1.18")
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	w := NewWebFingerprint()
	serverHeader := w.GetServer(server.URL)
	if serverHeader != "nginx/1.18" {
		t.Fatalf("expected nginx/1.18, got %q", serverHeader)
	}
}

func TestWebFingerprintReturnsEmptyServerOnError(t *testing.T) {
	w := NewWebFingerprint()
	result := w.GetServer("http://192.0.2.1:9999/nonexistent")
	if result != "" {
		t.Fatalf("expected empty string on error, got %q", result)
	}
}
