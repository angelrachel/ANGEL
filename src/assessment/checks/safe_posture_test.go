package checks

import "testing"

func TestSafePostureChecksAreRegistered(t *testing.T) {
	for _, expected := range []string{"web-cookie-policy", "web-cors-posture", "web-cache-control-policy", "web-information-disclosure", "web-server-header-disclosure", "recovery-backup-posture"} {
		if _, ok := Find(expected); !ok {
			t.Fatalf("check %s is not registered", expected)
		}
	}
}

func TestCookieAndCORSChecksRejectUnsafePolicies(t *testing.T) {
	cookie := CookiePolicy{}.Evaluate(Request{Headers: map[string]string{"Set-Cookie": "session=abc"}})
	if cookie.Passed || cookie.Severity != Medium {
		t.Fatalf("unexpected cookie result: %+v", cookie)
	}
	cors := CORSPosture{}.Evaluate(Request{Headers: map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Credentials": "true"}})
	if cors.Passed || cors.Severity != High {
		t.Fatalf("unexpected cors result: %+v", cors)
	}
}

func TestSafePostureChecksPassWithDefensiveConfiguration(t *testing.T) {
	checks := []struct {
		name  string
		value Result
	}{
		{"cookie", CookiePolicy{}.Evaluate(Request{Headers: map[string]string{"Set-Cookie": "session=abc; Secure; HttpOnly; SameSite=Lax"}})},
		{"cors", CORSPosture{}.Evaluate(Request{Headers: map[string]string{"Access-Control-Allow-Origin": "https://app.example", "Access-Control-Allow-Credentials": "true"}})},
		{"cache", CacheControlPolicy{}.Evaluate(Request{Headers: map[string]string{"Cache-Control": "private, no-store"}, Metadata: map[string]string{"sensitive_response": "true"}})},
		{"error", InformationDisclosure{}.Evaluate(Request{Body: "generic error"})},
		{"header", HeaderDisclosure{}.Evaluate(Request{Headers: map[string]string{}})},
		{"backup", BackupPosture{}.Evaluate(Request{Metadata: map[string]string{"backup_verified": "true"}})},
	}
	for _, item := range checks {
		if !item.value.Passed {
			t.Errorf("%s did not pass: %+v", item.name, item.value)
		}
	}
}
