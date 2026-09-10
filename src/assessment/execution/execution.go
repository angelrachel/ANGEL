package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ANGEL/src/assessment/checks"
	"ANGEL/src/assessment/collectors"
)

type Input struct {
	JobID   string         `json:"job_id"`
	Target  string         `json:"target"`
	CheckID string         `json:"check_id"`
	Fixture map[string]any `json:"fixture,omitempty"`
}
type Output struct {
	JobID       string                 `json:"job_id"`
	CheckID     string                 `json:"check_id"`
	Status      string                 `json:"status"`
	Result      checks.Result          `json:"result"`
	Observation collectors.Observation `json:"observation"`
	Evidence    []byte                 `json:"evidence"`
	CollectedAt time.Time              `json:"collected_at"`
}
type Adapter interface {
	Name() string
	Supports(checkID string) bool
	Run(context.Context, Input) (Output, error)
}
type Registry struct{ adapters []Adapter }

func NewRegistry() Registry {
	return Registry{adapters: []Adapter{FixtureAdapter{}, HTTPMetadataAdapter{}}}
}
func (r Registry) Resolve(checkID string) (Adapter, error) {
	if _, ok := checks.Find(checkID); !ok {
		return nil, fmt.Errorf("no safe adapter for check %s", checkID)
	}
	for _, adapter := range r.adapters {
		if adapter.Supports(checkID) {
			return adapter, nil
		}
	}
	return nil, fmt.Errorf("no safe adapter for check %s", checkID)
}
func (r Registry) Run(ctx context.Context, input Input) (Output, error) {
	adapter, err := r.Resolve(input.CheckID)
	if err != nil {
		return Output{}, err
	}
	return adapter.Run(ctx, input)
}

type FixtureAdapter struct{}

func (FixtureAdapter) Name() string                 { return "fixture-replay" }
func (FixtureAdapter) Supports(checkID string) bool { return checkID != "" }
func (FixtureAdapter) Run(_ context.Context, input Input) (Output, error) {
	if !strings.HasPrefix(input.Target, "fixture://") {
		return Output{}, fmt.Errorf("fixture adapter requires fixture target")
	}
	check, ok := checks.Find(input.CheckID)
	if !ok {
		return Output{}, fmt.Errorf("check not found")
	}
	headers := map[string]string{}
	if raw, ok := input.Fixture["headers"].(map[string]string); ok {
		headers = raw
	}
	request := checks.Request{ID: input.JobID, Target: input.Target, Headers: headers, Metadata: map[string]string{}}
	if value, ok := input.Fixture["body"].(string); ok {
		request.Body = value
	}
	if value, ok := input.Fixture["status_code"].(int); ok {
		request.StatusCode = value
	}
	result := check.Evaluate(request)
	evidence, _ := json.Marshal(map[string]any{"request": request, "result": result})
	return Output{JobID: input.JobID, CheckID: input.CheckID, Status: "COMPLETED", Result: result, Evidence: evidence, CollectedAt: time.Now().UTC()}, nil
}

type HTTPMetadataAdapter struct{ Client *http.Client }

func (HTTPMetadataAdapter) Name() string { return "http-metadata" }
func (HTTPMetadataAdapter) Supports(checkID string) bool {
	return checkID == "web-security-headers" || checkID == "web-method-policy" || checkID == "web-content-type-policy"
}
func (a HTTPMetadataAdapter) Run(ctx context.Context, input Input) (Output, error) {
	parsed, err := url.Parse(input.Target)
	if err != nil || parsed.Hostname() == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
		return Output{}, fmt.Errorf("HTTP adapter requires valid HTTP(S) target")
	}
	client := a.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodHead, input.Target, nil)
	if err != nil {
		return Output{}, err
	}
	response, err := client.Do(request)
	if err != nil {
		return Output{}, err
	}
	defer response.Body.Close()
	headers := map[string]string{}
	for _, key := range []string{"Content-Security-Policy", "X-Content-Type-Options", "Referrer-Policy", "Strict-Transport-Security", "Content-Type"} {
		headers[key] = response.Header.Get(key)
	}
	check, _ := checks.Find(input.CheckID)
	result := check.Evaluate(checks.Request{ID: input.JobID, Target: input.Target, StatusCode: response.StatusCode, Headers: headers})
	evidence, _ := json.Marshal(map[string]any{"status": response.StatusCode, "headers": headers})
	return Output{JobID: input.JobID, CheckID: input.CheckID, Status: "COMPLETED", Result: result, Evidence: evidence, CollectedAt: time.Now().UTC()}, nil
}
