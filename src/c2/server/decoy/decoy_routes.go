package decoy

import (
"net/http"
)

func SetDecoyHeaders(w http.ResponseWriter) {
w.Header().Set("Server", "nginx/1.24.0")
w.Header().Set("X-Content-Type-Options", "nosniff")
}

func ServeDecoyPage(w http.ResponseWriter, r *http.Request) {
SetDecoyHeaders(w)
ServeDecoy(w, r)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
SetDecoyHeaders(w)
w.WriteHeader(http.StatusNotFound)
w.Write([]byte("404 Not Found"))
}

func HandleOperator(w http.ResponseWriter, r *http.Request) {
SetDecoyHeaders(w)
w.WriteHeader(http.StatusOK)
w.Write([]byte("Operator Dashboard"))
}

func HandleBeacon(w http.ResponseWriter, r *http.Request) {
SetDecoyHeaders(w)
w.WriteHeader(http.StatusOK)
w.Write([]byte("Beacon Ready"))
}
