package assessmentapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"ANGEL/src/assessment/control"
	"ANGEL/src/assessment/plugins"
)

type Handler struct{ Controller *control.Controller }

func (h Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})
	mux.HandleFunc("/api/v1/modules", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, http.StatusOK, plugins.Catalog())
	})
	mux.HandleFunc("/api/v1/jobs", h.jobs)
	return mux
}
func (h Handler) jobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var request struct {
		EngagementID string `json:"engagement_id"`
		TaskType     string `json:"task_type"`
		Target       string `json:"target"`
		Action       string `json:"action"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	if strings.TrimSpace(request.EngagementID) == "" || strings.TrimSpace(request.Target) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "engagement_id and target are required"})
		return
	}
	job, err := h.Controller.Create(request.EngagementID, request.TaskType, request.Target, request.Action, time.Now().UTC())
	if err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusAccepted, job)
}
func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
