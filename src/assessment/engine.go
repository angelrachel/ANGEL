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

	"ANGEL/src/assessment/checks"
	"ANGEL/src/assessment/control"
	"ANGEL/src/assessment/domain"
	"ANGEL/src/assessment/evidence"
	"ANGEL/src/assessment/execution"
	"ANGEL/src/assessment/governance"
	"ANGEL/src/assessment/plugins"
	"ANGEL/src/assessment/policy"
	"ANGEL/src/assessment/ratelimit"
	"ANGEL/src/assessment/reporting"
	"ANGEL/src/assessment/risk"
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
	mu             sync.RWMutex
	agents         map[string]*Agent
	tasks          map[string]*Task
	results        map[string]*Result
	apiKey         string
	host           string
	port           int
	httpSrv        *http.Server
	policyEngine   *policy.Engine
	jobController  *control.Controller
	execution      execution.Registry
	evidenceSigner evidence.Signer
	evidence       map[string]evidence.Bundle
	findings       map[string]reporting.Finding
	reports        map[string]reporting.Report
	audit          *governance.AuditLog
	remediation    *risk.Tracker
	limiter        *ratelimit.Limiter
}

func NewEngine() *Engine {
	policyEngine := policy.NewEngine()
	jobController, err := control.NewController(policyEngine)
	if err != nil {
		panic(err)
	}
	signer, err := evidence.NewSigner()
	if err != nil {
		panic(err)
	}
	return &Engine{
		agents:         make(map[string]*Agent),
		tasks:          make(map[string]*Task),
		results:        make(map[string]*Result),
		apiKey:         os.Getenv("ANGEL_OPERATOR_KEY"),
		host:           envOr("ANGEL_HOST", "127.0.0.1"),
		port:           envPort("ANGEL_PORT", 8001),
		policyEngine:   policyEngine,
		jobController:  jobController,
		execution:      execution.NewRegistry(),
		evidenceSigner: signer,
		evidence:       make(map[string]evidence.Bundle),
		findings:       make(map[string]reporting.Finding),
		reports:        make(map[string]reporting.Report),
		audit:          governance.NewAuditLog(),
		remediation:    risk.NewTracker(),
		limiter:        ratelimit.New(20, 40),
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
	if e.limiter != nil {
		if err := e.limiter.Allow(clientIP(r.RemoteAddr), time.Now().UTC()); err != nil {
			return err
		}
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
	mux.HandleFunc("/api/v1/modules", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		writeJSON(w, http.StatusOK, plugins.Catalog())
	})
	mux.HandleFunc("/api/v1/checks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := e.auth(r); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			var request struct {
				CheckID string         `json:"check_id"`
				Input   checks.Request `json:"input"`
			}
			if err := decodeJSON(w, r, &request); err != nil {
				http.Error(w, "invalid check evaluation contract", http.StatusBadRequest)
				return
			}
			check, ok := checks.Find(request.CheckID)
			if !ok {
				http.Error(w, "check not found", http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, check.Evaluate(request.Input))
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		writeJSON(w, http.StatusOK, checks.IDs())
	})
	mux.HandleFunc("/api/v1/engagements", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var request struct {
			Engagement domain.Engagement   `json:"engagement"`
			Scope      []domain.ScopeEntry `json:"scope"`
			Budget     int                 `json:"budget"`
		}
		if err := decodeJSON(w, r, &request); err != nil {
			http.Error(w, "invalid engagement contract", http.StatusBadRequest)
			return
		}
		if err := e.policyEngine.RegisterEngagement(request.Engagement, request.Scope, request.Budget); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusCreated, request.Engagement)
	})
	mux.HandleFunc("/api/v1/assessment-jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var request struct {
			EngagementID string `json:"engagement_id"`
			TaskType     string `json:"task_type"`
			Target       string `json:"target"`
			Action       string `json:"action"`
		}
		if err := decodeJSON(w, r, &request); err != nil {
			http.Error(w, "invalid assessment job contract", http.StatusBadRequest)
			return
		}
		job, err := e.jobController.Create(request.EngagementID, request.TaskType, request.Target, request.Action, time.Now().UTC())
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		writeJSON(w, http.StatusAccepted, job)
	})
	mux.HandleFunc("/api/v1/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		writeJSON(w, http.StatusOK, e.jobController.List())
	})
	mux.HandleFunc("/api/v1/executions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var input execution.Input
		if err := decodeJSON(w, r, &input); err != nil || strings.TrimSpace(input.JobID) == "" || strings.TrimSpace(input.CheckID) == "" {
			http.Error(w, "invalid execution contract", http.StatusBadRequest)
			return
		}
		job, ok := e.jobController.Get(input.JobID)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}
		if job.Target != input.Target {
			http.Error(w, "execution target does not match signed job", http.StatusForbidden)
			return
		}
		output, err := e.execution.Run(r.Context(), input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		bundle, err := evidence.CreateBundle(job.EngagementID, job.ID, "application/json", output.Evidence, "0", true, output.CollectedAt, e.evidenceSigner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		finding := reporting.FromCheck(output.Result, input.Target, bundle.ID)
		e.mu.Lock()
		e.evidence[bundle.ID] = bundle
		e.findings[finding.ID] = finding
		e.mu.Unlock()
		_ = e.jobController.Transition(job.ID, control.Approved, output.CollectedAt)
		_ = e.jobController.Transition(job.ID, control.Queued, output.CollectedAt)
		_ = e.jobController.Transition(job.ID, control.Running, output.CollectedAt)
		_ = e.jobController.Transition(job.ID, control.Completed, output.CollectedAt)
		e.audit.Append("operator", "EXECUTION_COMPLETED", "job", job.ID, output.CollectedAt)
		writeJSON(w, http.StatusCreated, map[string]any{"output": output, "evidence": bundle, "finding": finding})
	})
	mux.HandleFunc("/api/v1/findings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		e.mu.RLock()
		items := make([]reporting.Finding, 0, len(e.findings))
		for _, item := range e.findings {
			items = append(items, item)
		}
		e.mu.RUnlock()
		writeJSON(w, http.StatusOK, items)
	})
	mux.HandleFunc("/api/v1/reports", func(w http.ResponseWriter, r *http.Request) {
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if r.Method == http.MethodGet {
			e.mu.RLock()
			items := make([]reporting.Report, 0, len(e.reports))
			for _, item := range e.reports {
				items = append(items, item)
			}
			e.mu.RUnlock()
			writeJSON(w, http.StatusOK, items)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var request struct {
			EngagementID string   `json:"engagement_id"`
			FindingIDs   []string `json:"finding_ids"`
		}
		if err := decodeJSON(w, r, &request); err != nil || strings.TrimSpace(request.EngagementID) == "" {
			http.Error(w, "invalid report contract", http.StatusBadRequest)
			return
		}
		e.mu.RLock()
		findings := make([]reporting.Finding, 0, len(request.FindingIDs))
		for _, id := range request.FindingIDs {
			if item, ok := e.findings[id]; ok {
				findings = append(findings, item)
			}
		}
		e.mu.RUnlock()
		report, err := reporting.New(request.EngagementID, findings, time.Now().UTC())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		e.mu.Lock()
		e.reports[report.ID] = report
		e.mu.Unlock()
		e.audit.Append("operator", "REPORT_GENERATED", "report", report.ID, report.GeneratedAt)
		writeJSON(w, http.StatusCreated, report)
	})
	mux.HandleFunc("/api/v1/evidence", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		e.mu.RLock()
		items := make([]evidence.Bundle, 0, len(e.evidence))
		for _, item := range e.evidence {
			safe := item
			safe.Payload = nil
			items = append(items, safe)
		}
		e.mu.RUnlock()
		writeJSON(w, http.StatusOK, items)
	})
	mux.HandleFunc("/api/v1/remediations", func(w http.ResponseWriter, r *http.Request) {
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, e.remediation.List())
			return
		}
		if r.Method == http.MethodPost {
			var request struct {
				FindingID string    `json:"finding_id"`
				Owner     string    `json:"owner"`
				Plan      string    `json:"plan"`
				DueAt     time.Time `json:"due_at"`
			}
			if err := decodeJSON(w, r, &request); err != nil {
				http.Error(w, "invalid remediation contract", http.StatusBadRequest)
				return
			}
			item, err := e.remediation.Create(request.FindingID, request.Owner, request.Plan, request.DueAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			e.audit.Append("operator", "REMEDIATION_CREATED", "remediation", item.ID, time.Now().UTC())
			writeJSON(w, http.StatusCreated, item)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/api/v1/audit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := e.auth(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		writeJSON(w, http.StatusOK, e.audit.List())
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
