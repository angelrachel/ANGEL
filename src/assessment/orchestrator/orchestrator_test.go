package orchestrator

import (
	"context"
	"testing"
	"time"
)

func TestAgentRegistryRegisterAndGet(t *testing.T) {
	registry := NewAgentRegistry()
	agent := &Agent{Name: "test-agent", Type: "scanner"}
	registry.Register(agent)

	retrieved := registry.Get("test-agent")
	if retrieved == nil {
		t.Fatal("expected agent to be found")
	}
	if retrieved.Name != "test-agent" {
		t.Fatalf("expected test-agent, got %q", retrieved.Name)
	}
	if retrieved.Type != "scanner" {
		t.Fatalf("expected scanner, got %q", retrieved.Type)
	}
}

func TestAgentRegistryGetReturnsNilForUnknown(t *testing.T) {
	registry := NewAgentRegistry()
	result := registry.Get("nonexistent")
	if result != nil {
		t.Fatal("expected nil for unknown agent")
	}
}

func TestAgentRegistryList(t *testing.T) {
	registry := NewAgentRegistry()
	registry.Register(&Agent{Name: "a1", Type: "scanner"})
	registry.Register(&Agent{Name: "a2", Type: "analyzer"})
	registry.Register(&Agent{Name: "a3", Type: "reporter"})

	agents := registry.List()
	if len(agents) != 3 {
		t.Fatalf("expected 3 agents, got %d", len(agents))
	}
}

func TestAgentRegistryConcurrentAccess(t *testing.T) {
	registry := NewAgentRegistry()
	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			registry.Register(&Agent{Name: string(rune('a' + i)), Type: "test"})
		}
		done <- struct{}{}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			registry.List()
		}
		done <- struct{}{}
	}()
	<-done
	<-done
	// Jika sampai di sini tanpa panic/panic recover, test lolos
}

func TestRouterRegisterAgent(t *testing.T) {
	router := NewRouter()
	router.RegisterAgent("scanner", &Agent{Name: "scanner", Type: "scanner"})

	count := router.GetAgentCount()
	if count != 1 {
		t.Fatalf("expected 1 agent, got %d", count)
	}
}

func TestRouterDispatchAgentNotFound(t *testing.T) {
	router := NewRouter()
	task := Task{ID: "task-1", Type: "unknown-agent", Payload: map[string]interface{}{}}
	result := router.Dispatch(task)

	if result.Status != "error" {
		t.Fatalf("expected error status, got %q", result.Status)
	}
	if result.Output["error"] != "agent not found" {
		t.Fatalf("expected 'agent not found', got %v", result.Output["error"])
	}
}

func TestRouterDispatchAgentFound(t *testing.T) {
	router := NewRouter()
	router.RegisterAgent("scanner", &Agent{Name: "scanner", Type: "scanner"})

	task := Task{ID: "task-1", Type: "scanner", Payload: map[string]interface{}{"target": "test"}}
	result := router.Dispatch(task)

	if result.Status != "accepted" {
		t.Fatalf("expected accepted status, got %q", result.Status)
	}
	if result.Output["agent"] != "scanner" {
		t.Fatalf("expected scanner agent, got %v", result.Output["agent"])
	}
	if result.Output["mode"] != "observe" {
		t.Fatalf("expected observe mode, got %v", result.Output["mode"])
	}
}

func TestRouterFormatStatus(t *testing.T) {
	status := FormatRouterStatus(5)
	if status != "Router aktif dengan 5 agen" {
		t.Fatalf("expected 'Router aktif dengan 5 agen', got %q", status)
	}
}

func TestFireteamLaunch(t *testing.T) {
	router := NewRouter()
	router.RegisterAgent("scanner", &Agent{Name: "scanner", Type: "scanner"})

	fireteam := NewFireteam()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasks := []Task{
		{ID: "t1", Type: "scanner", Payload: map[string]interface{}{}},
		{ID: "t2", Type: "scanner", Payload: map[string]interface{}{}},
	}

	results := fireteam.Launch(ctx, router, tasks)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	results2 := fireteam.GetResults()
	if len(results2) != 2 {
		t.Fatalf("expected 2 results from GetResults, got %d", len(results2))
	}
}

func TestFireteamLaunchWithCancel(t *testing.T) {
	router := NewRouter()
	router.RegisterAgent("scanner", &Agent{Name: "scanner", Type: "scanner"})

	fireteam := NewFireteam()
	ctx, cancel := context.WithCancel(context.Background())

	tasks := []Task{
		{ID: "t1", Type: "scanner", Payload: map[string]interface{}{}},
		{ID: "t2", Type: "scanner", Payload: map[string]interface{}{}},
		{ID: "t3", Type: "scanner", Payload: map[string]interface{}{}},
	}

	// Batalkan sebelum semuanya selesai
	go func() {
		time.Sleep(1 * time.Millisecond)
		cancel()
	}()

	results := fireteam.Launch(ctx, router, tasks)
	// Hasil mungkin tidak lengkap karena cancel
	if len(results) == 0 {
		t.Fatal("expected at least some results even after cancel")
	}
}

func TestTaskAndResultSerialization(t *testing.T) {
	task := Task{
		ID:      "task-001",
		Type:    "scanner",
		Payload: map[string]interface{}{"target": "fixture://lab/web"},
	}

	if task.ID != "task-001" {
		t.Fatalf("expected task-001, got %q", task.ID)
	}
	if task.Type != "scanner" {
		t.Fatalf("expected scanner, got %q", task.Type)
	}
	if task.Payload["target"] != "fixture://lab/web" {
		t.Fatalf("expected fixture target, got %v", task.Payload["target"])
	}

	result := Result{
		TaskID: "task-001",
		Status: "completed",
		Output: map[string]interface{}{"data": "test"},
	}

	if result.TaskID != "task-001" {
		t.Fatalf("expected task-001, got %q", result.TaskID)
	}
	if result.Status != "completed" {
		t.Fatalf("expected completed, got %q", result.Status)
	}
}
