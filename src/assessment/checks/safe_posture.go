package checks

import (
	"fmt"
	"strings"
)

type CookiePolicy struct{}

func (CookiePolicy) ID() string { return "web-cookie-policy" }
func (CookiePolicy) Layer() int { return 9 }
func (CookiePolicy) Evaluate(r Request) Result {
	value := strings.ToLower(r.Headers["Set-Cookie"])
	required := []string{"secure", "httponly", "samesite"}
	missing := make([]string, 0, len(required))
	for _, flag := range required {
		if !strings.Contains(value, flag) {
			missing = append(missing, flag)
		}
	}
	return result("web-cookie-policy", 9, "Cookie security policy", len(missing) == 0, Medium, .94, fmt.Sprintf("missing cookie attributes: %s", strings.Join(missing, ", ")), map[string]string{"missing": strings.Join(missing, ",")}, "Set Secure, HttpOnly, and SameSite attributes on session cookies")
}

type CORSPosture struct{}

func (CORSPosture) ID() string { return "web-cors-posture" }
func (CORSPosture) Layer() int { return 9 }
func (CORSPosture) Evaluate(r Request) Result {
	origin := strings.TrimSpace(r.Headers["Access-Control-Allow-Origin"])
	credentials := strings.EqualFold(strings.TrimSpace(r.Headers["Access-Control-Allow-Credentials"]), "true")
	unsafe := origin == "*" && credentials
	return result("web-cors-posture", 9, "CORS posture", !unsafe, High, .95, "CORS response policy evaluated", map[string]string{"origin": origin, "credentials": fmt.Sprint(credentials)}, "Do not combine wildcard origins with credentialed cross-origin requests")
}

type CacheControlPolicy struct{}

func (CacheControlPolicy) ID() string { return "web-cache-control-policy" }
func (CacheControlPolicy) Layer() int { return 14 }
func (CacheControlPolicy) Evaluate(r Request) Result {
	cache := strings.ToLower(r.Headers["Cache-Control"])
	sensitive := strings.EqualFold(r.Metadata["sensitive_response"], "true")
	passed := !sensitive || (strings.Contains(cache, "no-store") && strings.Contains(cache, "private"))
	return result("web-cache-control-policy", 14, "Sensitive response cache policy", passed, High, .93, "cache-control policy evaluated", map[string]string{"cache_control": cache, "sensitive": fmt.Sprint(sensitive)}, "Use private, no-store caching for sensitive responses")
}

type InformationDisclosure struct{}

func (InformationDisclosure) ID() string { return "web-information-disclosure" }
func (InformationDisclosure) Layer() int { return 14 }
func (InformationDisclosure) Evaluate(r Request) Result {
	body := strings.ToLower(r.Body)
	markers := []string{"stack trace", "debug=true", "sqlstate", "at /app/", "traceback"}
	found := ""
	for _, marker := range markers {
		if strings.Contains(body, marker) {
			found = marker
			break
		}
	}
	return result("web-information-disclosure", 14, "Error information disclosure", found == "", Medium, .9, "error response disclosure scan", map[string]string{"marker": found}, "Return generic errors externally and keep diagnostic details in protected logs")
}

type HeaderDisclosure struct{}

func (HeaderDisclosure) ID() string { return "web-server-header-disclosure" }
func (HeaderDisclosure) Layer() int { return 9 }
func (HeaderDisclosure) Evaluate(r Request) Result {
	server := strings.TrimSpace(r.Headers["Server"])
	powered := strings.TrimSpace(r.Headers["X-Powered-By"])
	passed := server == "" && powered == ""
	return result("web-server-header-disclosure", 9, "Server header disclosure", passed, Low, .91, "server identification headers evaluated", map[string]string{"server": server, "powered_by": powered}, "Remove unnecessary implementation-identifying response headers")
}

type BackupPosture struct{}

func (BackupPosture) ID() string { return "recovery-backup-posture" }
func (BackupPosture) Layer() int { return 21 }
func (BackupPosture) Evaluate(r Request) Result {
	verified := strings.EqualFold(r.Metadata["backup_verified"], "true")
	return result("recovery-backup-posture", 21, "Backup verification posture", verified, High, .9, "backup restore evidence evaluated", map[string]string{"backup_verified": r.Metadata["backup_verified"]}, "Verify backup integrity, retention, restoration, and ownership through a controlled test")
}
