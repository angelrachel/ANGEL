package decoy

import (
"net/http"
)

func ServeDecoy(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "text/html; charset=utf-8")
w.Header().Set("Server", "nginx/1.24.0")
w.Header().Set("X-Content-Type-Options", "nosniff")
w.WriteHeader(http.StatusOK)
w.Write([]byte(DecoyHTML))
}
