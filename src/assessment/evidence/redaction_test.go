package evidence

import (
	"strings"
	"testing"
)

func TestRedactorRemovesSensitiveValues(t *testing.T) {
	input := `password="secret" api_key='key' secret="value" Authorization: Bearer abc.def-123 123-45-6789`
	output := NewRedactor().Redact(input)
	if strings.Contains(output, "secret") || strings.Contains(output, "abc.def") || strings.Contains(output, "123-45-6789") {
		t.Fatalf("sensitive value remained after redaction: %q", output)
	}
	if strings.Count(output, "[REDACTED]") < 5 {
		t.Fatalf("expected all sensitive patterns to be redacted: %q", output)
	}
}

func TestNilRedactorLeavesInputUnchanged(t *testing.T) {
	input := "synthetic evidence"
	var redactor *Redactor
	if got := redactor.Redact(input); got != input {
		t.Fatalf("expected nil redactor to preserve input, got %q", got)
	}
}
