package listener

import (
"context"
"net/http"
"time"
)

type HTTPServer struct {
Server *http.Server
Port   int
}

func NewHTTPServer(port int) *HTTPServer {
mux := http.NewServeMux()
return &HTTPServer{
Server: &http.Server{
Addr:              ":" + itoa(port),
Handler:           mux,
ReadTimeout:       10 * time.Second,
WriteTimeout:      10 * time.Second,
IdleTimeout:       60 * time.Second,
ReadHeaderTimeout: 5 * time.Second,
},
Port: port,
}
}

func (h *HTTPServer) ListenAndServe(ctx context.Context) error {
go func() {
<-ctx.Done()
h.Server.Shutdown(context.Background())
}()
return h.Server.ListenAndServe()
}

func (h *HTTPServer) RegisterHandler(pattern string, handler http.HandlerFunc) {
h.Server.Handler.(*http.ServeMux).HandleFunc(pattern, handler)
}

func itoa(i int) string {
if i == 0 {
return "0"
}
digits := []byte{}
for i > 0 {
digits = append([]byte{byte('0' + i%10)}, digits...)
i = i / 10
}
return string(digits)
}
