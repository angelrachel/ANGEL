package main

import (
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"encoding/base64"
"encoding/json"
"fmt"
"io"
"log"
"net/http"
"strings"
"sync"
"time"
)

type Agent struct {
ID       string `json:"id"`
Hostname string `json:"hostname"`
OS       string `json:"os"`
LastSeen int64  `json:"last_seen"`
}

type Task struct {
ID        string `json:"id"`
AgentID   string `json:"agent_id"`
Command   string `json:"command"`
Status    string `json:"status"`
CreatedAt int64  `json:"created_at"`
}

var (
agents   = make(map[string]Agent)
tasks    = make(map[string][]Task)
taskID   = 0
agentsMu sync.Mutex
tasksMu  sync.Mutex
)

const API_TOKEN = "my-super-secret-token-123"

var tokenRoles = map[string]string{
"my-super-secret-token-123": "admin",
"operator-token-456":        "operator",
"viewer-token-789":          "viewer",
}

var allowedCommands = map[string]bool{
"whoami": true, "hostname": true, "id": true, "ls": true,
"pwd": true, "date": true, "uptime": true, "echo": true,
"ping": true, "ipconfig": true, "ifconfig": true, "netstat": true,
"ps": true, "top": true, "dir": true,
}

func getRole(r *http.Request) string {
authHeader := r.Header.Get("Authorization")
if authHeader == "" {
return ""
}
parts := strings.Split(authHeader, " ")
if len(parts) != 2 || parts[0] != "Bearer" {
return ""
}
return tokenRoles[parts[1]]
}

func isCommandAllowed(cmd string) bool {
parts := strings.Fields(cmd)
if len(parts) == 0 {
return false
}
return allowedCommands[parts[0]]
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
return func(w http.ResponseWriter, r *http.Request) {
if getRole(r) == "" {
http.Error(w, "Unauthorized", http.StatusUnauthorized)
return
}
next(w, r)
}
}

func XOREncrypt(key, plaintext []byte) (string, error) {
if len(key) == 0 {
return "", fmt.Errorf("key cannot be empty")
}
ciphertext := make([]byte, len(plaintext))
for i := 0; i < len(plaintext); i++ {
ciphertext[i] = plaintext[i] ^ key[i%len(key)]
}
return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func XORDecrypt(key []byte, cryptoText string) ([]byte, error) {
if len(key) == 0 {
return nil, fmt.Errorf("key cannot be empty")
}
ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
if err != nil {
return nil, err
}
plaintext := make([]byte, len(ciphertext))
for i := 0; i < len(ciphertext); i++ {
plaintext[i] = ciphertext[i] ^ key[i%len(key)]
}
return plaintext, nil
}

func AESGCMEncrypt(key, plaintext []byte) (string, error) {
block, err := aes.NewCipher(key)
if err != nil {
return "", err
}
gcm, err := cipher.NewGCM(block)
if err != nil {
return "", err
}
nonce := make([]byte, gcm.NonceSize())
if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
return "", err
}
ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func AESGCMDecrypt(key []byte, cryptoText string) ([]byte, error) {
ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
if err != nil {
return nil, err
}
block, err := aes.NewCipher(key)
if err != nil {
return nil, err
}
gcm, err := cipher.NewGCM(block)
if err != nil {
return nil, err
}
if len(ciphertext) < gcm.NonceSize() {
return nil, fmt.Errorf("ciphertext too short")
}
nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
return gcm.Open(nil, nonce, ciphertext, nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}
var agent Agent
if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
http.Error(w, "Invalid request", http.StatusBadRequest)
return
}
agentsMu.Lock()
agent.LastSeen = time.Now().Unix()
agents[agent.ID] = agent
agentsMu.Unlock()
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}
var req map[string]string
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request", http.StatusBadRequest)
return
}
agentID := req["agent_id"]
agentsMu.Lock()
_, exists := agents[agentID]
agentsMu.Unlock()
if !exists {
http.Error(w, "Agent not found", http.StatusNotFound)
return
}
if cmd, ok := req["command"]; ok && cmd != "" {
role := getRole(r)
if role != "admin" && role != "operator" {
http.Error(w, "Forbidden", http.StatusForbidden)
return
}
if !isCommandAllowed(cmd) {
http.Error(w, fmt.Sprintf("Command not allowed: %s", cmd), http.StatusForbidden)
return
}
tasksMu.Lock()
taskID++
task := Task{ID: fmt.Sprintf("task-%d", taskID), AgentID: agentID, Command: cmd, Status: "pending", CreatedAt: time.Now().Unix()}
tasks[agentID] = append(tasks[agentID], task)
tasksMu.Unlock()
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]string{"status": "task_added"})
return
}
tasksMu.Lock()
agentTasks := tasks[agentID]
var pendingTask *Task
for i := range agentTasks {
if agentTasks[i].Status == "pending" {
pendingTask = &agentTasks[i]
break
}
}
tasksMu.Unlock()
if pendingTask != nil {
tasksMu.Lock()
for i := range tasks[agentID] {
if tasks[agentID][i].ID == pendingTask.ID {
tasks[agentID][i].Status = "assigned"
break
}
}
tasksMu.Unlock()
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]interface{}{"command": pendingTask.Command})
} else {
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]interface{}{"command": nil})
}
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}
var result map[string]string
if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
http.Error(w, "Invalid request", http.StatusBadRequest)
return
}
fmt.Printf("[+] Result from %s: %s\n", result["agent_id"], result["output"])
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]string{"status": "recorded"})
}

func engineStatusHandler(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]interface{}{
"status":    "active",
"time":      time.Now().Unix(),
"version":   "ANGEL-P0-ENGINE",
"rbac":      "active",
"allowlist": "active",
})
}

func main() {
http.HandleFunc("/register", authMiddleware(registerHandler))
http.HandleFunc("/task", authMiddleware(taskHandler))
http.HandleFunc("/result", authMiddleware(resultHandler))
http.HandleFunc("/status", engineStatusHandler)

fmt.Println("[+] ANGEL P0 Platform Engine Running on port 8001")
log.Fatal(http.ListenAndServe(":8001", nil))
}
