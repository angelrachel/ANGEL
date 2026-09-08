package osint

import (
"io"
"net/http"
"strings"
"time"
)

type WebFingerprint struct {
Client *http.Client
}

func NewWebFingerprint() *WebFingerprint {
return &WebFingerprint{
Client: &http.Client{Timeout: 10 * time.Second},
}
}

func (w *WebFingerprint) Fingerprint(url string) string {
resp, err := w.Client.Get(url)
if err != nil {
return ""
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
content := string(body)
if strings.Contains(content, "wp-content") {
return "WordPress"
}
if strings.Contains(content, "Joomla") {
return "Joomla"
}
if strings.Contains(content, "Drupal") {
return "Drupal"
}
if strings.Contains(content, "Laravel") {
return "Laravel"
}
return "Unknown"
}

func (w *WebFingerprint) GetServer(url string) string {
resp, err := w.Client.Get(url)
if err != nil {
return ""
}
defer resp.Body.Close()
return resp.Header.Get("Server")
}

func (w *WebFingerprint) GetPoweredBy(url string) string {
resp, err := w.Client.Get(url)
if err != nil {
return ""
}
defer resp.Body.Close()
return resp.Header.Get("X-Powered-By")
}
