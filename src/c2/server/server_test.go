package main

import (
"bytes"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"
)

func TestRegisterHandler(t *testing.T) {
reqBody := `{"id":"test-agent","hostname":"test","os":"linux"}`
req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(reqBody))
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer my-super-secret-token-123")

w := httptest.NewRecorder()
registerHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status 200, got %d", w.Code)
}

var resp map[string]string
json.NewDecoder(w.Body).Decode(&resp)
if resp["status"] != "registered" {
t.Errorf("Expected status 'registered', got '%s'", resp["status"])
}
}

func TestTaskHandler(t *testing.T) {
// Register agent first
registerReq := `{"id":"test-agent","hostname":"test","os":"linux"}`
req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(registerReq))
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer my-super-secret-token-123")
w := httptest.NewRecorder()
registerHandler(w, req)

// Send task
taskReq := `{"agent_id":"test-agent","command":"whoami"}`
req = httptest.NewRequest("POST", "/task", bytes.NewBufferString(taskReq))
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer my-super-secret-token-123")
w = httptest.NewRecorder()
taskHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status 200, got %d", w.Code)
}

var resp map[string]string
json.NewDecoder(w.Body).Decode(&resp)
if resp["status"] != "task_added" {
t.Errorf("Expected status 'task_added', got '%s'", resp["status"])
}
}
