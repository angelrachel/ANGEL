package collectors

import (
	"context"
	"testing"
)

func TestDependencyManifestIsValidatedAndSorted(t *testing.T) {
	deps, err := ParseDependencyManifest([]byte(`[{"name":"z","version":"1"},{"name":"a","version":"2"}]`))
	if err != nil {
		t.Fatal(err)
	}
	if deps[0].Name != "a" {
		t.Fatalf("not sorted: %#v", deps)
	}
}
func TestAuthorizationMatrixProducesDeterministicObservations(t *testing.T) {
	result := EvaluateAuthorizationMatrix([]AuthorizationCase{{Subject: "user-b", Resource: "/a", Action: "read", Expected: "deny", Actual: 403}, {Subject: "user-a", Resource: "/a", Action: "read", Expected: "allow", Actual: 200}})
	if len(result) != 2 || result[0].Attributes["subject"] != "user-a" || result[1].Attributes["passed"] != "true" {
		t.Fatalf("unexpected matrix: %#v", result)
	}
}
func TestCloudPostureRejectsInvalidJSON(t *testing.T) {
	if _, err := ParseCloudPosture([]byte("bad")); err == nil {
		t.Fatal("invalid posture accepted")
	}
}
func TestHTTPCollectorRejectsUnsupportedTarget(t *testing.T) {
	_, err := (HTTPCollector{}).Collect(context.Background(), "ftp://example.test")
	if err == nil {
		t.Fatal("unsupported target accepted")
	}
}
