package main

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Agent struct {
	ID       string    `json:"id"`
	Hostname string    `json:"hostname"`
	OS       string    `json:"os"`
	IP       string    `json:"ip"`
	LastSeen time.Time `json:"last_seen"`
}

type Task struct {
	ID          string                 `json:"id"`
	AgentID     string                 `json:"agent_id"`
	Type        string                 `json:"type"`
	TargetRef   string                 `json:"target_ref"`
	Mode        string                 `json:"mode"`
	RequestedBy string                 `json:"requested_by"`
	Status      string                 `json:"status"`
	Payload     map[string]interface{} `json:"payload"`
}

type Result struct {
	TaskID  string                 `json:"task_id"`
	AgentID string                 `json:"agent_id"`
	Output  map[string]interface{} `json:"output"`
	Status  string                 `json:"status"`
}

type Engine struct {
	mu      sync.RWMutex
	agents  map[string]*Agent
	tasks   map[string]*Task
	results map[string]*Result
	apiKey  string
	host    string
	port    int
	httpSrv *http.Server
}

func NewEngine() *Engine {
	return &Engine{
		agents:  make(map[string]*Agent),
		tasks:   make(map[string]*Task),
		results: make(map[string]*Result),
		apiKey:  os.Getenv("ANGEL_OPERATOR_KEY"),
		host:    envOr("ANGEL_HOST", "127.0.0.1"),
		port:    envPort("ANGEL_PORT", 8001),
	}
}

func envOr(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func envPort(name string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(os.Getenv(name)))
	if err != nil || value < 1 || value > 65535 {
		return fallback
	}
	return value
}

func (e *Engine) address() string {
	return net.JoinHostPort(e.host, strconv.Itoa(e.port))
}

func (e *Engine) auth(r *http.Request) error {
	if e.apiKey == "" {
		return errors.New("operator authentication is not configured")
	}
	const prefix = "Bearer "
	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, prefix) {
		return errors.New("missing bearer token")
	}
	provided := strings.TrimSpace(strings.TrimPrefix(token, prefix))
	if provided == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(e.apiKey)) != 1 {
		return errors.New("invalid token")
	}
	return nil
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	var extra interface{}
	if err := decoder.Decode(&extra); err != io.EOF {
		return errors.New("request contains multiple JSON values")
	}
	return nil
}

func clientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return host
	}
	return remoteAddr
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func (e *Engine) withHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (e *Engine) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/api/v1/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var a Agent
		if err := decodeJSON(w, r, &a); err != nil || strings.TrimSpace(a.ID) == "" {
			http.Error(w, "bad json or missing agent id", http.StatusBadRequest)
			return
		}
		a.LastSeen = time.Now().UTC()
		a.IP = clientIP(r.RemoteAddr)
		e.mu.Lock()
		e.agents[a.ID] = &a
		e.mu.Unlock()
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if r.Method == http.MethodPost {
			var task Task
			if err := decodeJSON(w, r, &task); err != nil {
				http.Error(w, "invalid task contract", http.StatusBadRequest)
				return
			}
			if strings.TrimSpace(task.ID) == "" || strings.TrimSpace(task.AgentID) == "" || strings.TrimSpace(task.Type) == "" || strings.TrimSpace(task.RequestedBy) == "" {
				http.Error(w, "task identity fields are required", http.StatusBadRequest)
				return
			}
			if !strings.HasPrefix(strings.TrimSpace(task.TargetRef), "fixture://") || len(strings.TrimSpace(task.TargetRef)) <= len("fixture://") {
				http.Error(w, "target_ref must use a fixture reference", http.StatusForbidden)
				return
			}
			if task.Mode != "observe" && task.Mode != "simulate" {
				http.Error(w, "mode must be observe or simulate", http.StatusBadRequest)
				return
			}
			e.mu.Lock()
			if _, exists := e.agents[task.AgentID]; !exists {
				e.mu.Unlock()
				http.Error(w, "agent is not registered", http.StatusForbidden)
				return
			}
			if _, exists := e.tasks[task.ID]; exists {
				e.mu.Unlock()
				writeJSON(w, http.StatusOK, map[string]string{"status": "exists"})
				return
			}
			task.Status = "queued"
			e.tasks[task.ID] = &task
			e.mu.Unlock()
			writeJSON(w, http.StatusAccepted, task)
			return
		}

		agentID := strings.TrimSpace(r.URL.Query().Get("agent_id"))
		if agentID == "" {
			http.Error(w, "agent_id is required", http.StatusBadRequest)
			return
		}
		e.mu.RLock()
		tasks := make([]Task, 0)
		for _, t := range e.tasks {
			if t.AgentID == agentID {
				tasks = append(tasks, *t)
			}
		}
		e.mu.RUnlock()
		writeJSON(w, http.StatusOK, tasks)
	})

	mux.HandleFunc("/api/v1/results", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var res Result
		if err := decodeJSON(w, r, &res); err != nil || strings.TrimSpace(res.TaskID) == "" || strings.TrimSpace(res.AgentID) == "" {
			http.Error(w, "bad json or missing result identity", http.StatusBadRequest)
			return
		}
		e.mu.Lock()
		e.results[res.TaskID] = &res
		delete(e.tasks, res.TaskID)
		e.mu.Unlock()
		writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
	})

	e.httpSrv = &http.Server{
		Addr:              e.address(),
		Handler:           e.withHeaders(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = e.httpSrv.Shutdown(shutdownCtx)
	}()
	fmt.Printf("[+] ANGEL ASSESSMENT CONTROL PLANE RUNNING ON %s\n", e.address())
	if err := e.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	engine := NewEngine()
	if err := engine.Start(ctx); err != nil {
		fmt.Println("[FATAL] Engine failed:", err)
		os.Exit(1)
	}
}
