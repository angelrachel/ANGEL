package listener

import (
"context"
"encoding/json"
"net/http"
)

type WebSocketMessage struct {
Type    string                 `json:"type"`
Payload map[string]interface{} `json:"payload"`
}

type WebSocketHandler struct {
OnMessage func(message WebSocketMessage)
}

func (w *WebSocketHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
if w.OnMessage == nil {
return
}
var message WebSocketMessage
if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
return
}
w.OnMessage(message)
}

type WebSocketListener struct {
Port    int
Handler *WebSocketHandler
}

func NewWebSocketListener(port int, handler *WebSocketHandler) *WebSocketListener {
return &WebSocketListener{Port: port, Handler: handler}
}

func (w *WebSocketListener) Listen(ctx context.Context) error {
mux := http.NewServeMux()
mux.Handle("/ws", w.Handler)
server := &http.Server{Addr: ":" + itoa(w.Port), Handler: mux}
go func() {
<-ctx.Done()
server.Shutdown(context.Background())
}()
return server.ListenAndServe()
}

func (w *WebSocketListener) RegisterHandler(pattern string, handler http.Handler) {
http.Handle(pattern, handler)
}
