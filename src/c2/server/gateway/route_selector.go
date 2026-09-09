package gateway

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
)

type AuthError struct {
	Message string
}

func (e AuthError) Error() string {
	return e.Message
}

func IsBeaconRequest(r *http.Request) error {
	token := strings.TrimSpace(r.Header.Get("X-Beacon-Token"))
	configured := strings.TrimSpace(os.Getenv("ANGEL_BEACON_TOKEN"))
	if token == "" || configured == "" {
		return AuthError{Message: "beacon authentication is not configured"}
	}
	if subtle.ConstantTimeCompare([]byte(token), []byte(configured)) != 1 {
		return AuthError{Message: "invalid beacon token"}
	}
	return nil
}

func IsOperatorRequest(r *http.Request) error {
	key := strings.TrimSpace(r.Header.Get("X-Operator-Key"))
	configured := strings.TrimSpace(os.Getenv("ANGEL_OPERATOR_KEY"))
	if key == "" || configured == "" {
		return AuthError{Message: "operator authentication is not configured"}
	}
	if subtle.ConstantTimeCompare([]byte(key), []byte(configured)) != 1 {
		return AuthError{Message: "invalid operator key"}
	}
	return nil
}

func SelectRoute(r *http.Request) string {
	path := r.URL.Path
	if strings.HasPrefix(path, "/beacon/") && IsBeaconRequest(r) == nil {
		return "beacon"
	}
	if strings.HasPrefix(path, "/admin/") && IsOperatorRequest(r) == nil {
		return "operator"
	}
	return "decoy"
}
