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
token := r.Header.Get("X-Beacon-Token")
if token == "" {
return AuthError{Message: "missing beacon token"}
}
if subtle.ConstantTimeCompare([]byte(token), []byte(os.Getenv("ANGEL_BEACON_TOKEN"))) != 1 {
return AuthError{Message: "invalid beacon token"}
}
return AuthError{Message: "success"}
}

func IsOperatorRequest(r *http.Request) error {
key := r.Header.Get("X-Operator-Key")
if key == "" {
return AuthError{Message: "missing operator key"}
}
if subtle.ConstantTimeCompare([]byte(key), []byte(os.Getenv("ANGEL_OPERATOR_KEY"))) != 1 {
return AuthError{Message: "invalid operator key"}
}
return AuthError{Message: "success"}
}

func SelectRoute(r *http.Request) string {
path := r.URL.Path
if strings.HasPrefix(path, "/beacon/") && IsBeaconRequest(r) == (AuthError{}) {
return "beacon"
}
if strings.HasPrefix(path, "/admin/") && IsOperatorRequest(r) == (AuthError{}) {
return "operator"
}
return "decoy"
}
