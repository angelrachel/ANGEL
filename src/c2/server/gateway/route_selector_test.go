package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSelectRouteRequiresConfiguredToken(t *testing.T) {
	t.Setenv("ANGEL_BEACON_TOKEN", "beacon-secret")
	t.Setenv("ANGEL_OPERATOR_KEY", "operator-secret")

	beacon := httptest.NewRequest(http.MethodGet, "/beacon/poll", nil)
	beacon.Header.Set("X-Beacon-Token", "beacon-secret")
	if got := SelectRoute(beacon); got != "beacon" {
		t.Fatalf("beacon route = %q", got)
	}

	operator := httptest.NewRequest(http.MethodGet, "/admin/status", nil)
	operator.Header.Set("X-Operator-Key", "operator-secret")
	if got := SelectRoute(operator); got != "operator" {
		t.Fatalf("operator route = %q", got)
	}

	missing := httptest.NewRequest(http.MethodGet, "/admin/status", nil)
	if got := SelectRoute(missing); got != "decoy" {
		t.Fatalf("missing token route = %q", got)
	}
}

func TestSelectRouteFailsClosedWhenSecretUnset(t *testing.T) {
	t.Setenv("ANGEL_BEACON_TOKEN", "")
	t.Setenv("ANGEL_OPERATOR_KEY", "")
	request := httptest.NewRequest(http.MethodGet, "/beacon/poll", nil)
	request.Header.Set("X-Beacon-Token", "anything")
	if got := SelectRoute(request); got != "decoy" {
		t.Fatalf("unset secret route = %q", got)
	}
}
