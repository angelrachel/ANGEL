package decoy

import (
"errors"
"net/http"
)

func SelectRoute(r *http.Request) (string, error) {
if r.Header.Get("X-Beacon-Token") != "" {
return "beacon", errors.New("beacon request detected")
}
if r.Header.Get("X-Operator-Key") != "" {
return "operator", errors.New("operator request detected")
}
return "unknown", errors.New("unknown request type")
}

func ServeDecoyOrRedirect(w http.ResponseWriter, r *http.Request) {
route, err := SelectRoute(r)
if err != nil {
if route == "beacon" || route == "operator" {
http.Redirect(w, r, "/admin", http.StatusFound)
return
}
ServeDecoy(w, r)
}
}
