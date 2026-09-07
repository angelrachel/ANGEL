package osint

import (
"crypto/tls"
"fmt"
"net/http"
"regexp"
"strings"
"time"
)

type WebOSINT struct {
Client *http.Client
}

func NewWebOSINT() *WebOSINT {
return &WebOSINT{
Client: &http.Client{
Timeout: 10 * time.Second,
Transport: &http.Transport{
TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
},
},
}
}

func (w *WebOSINT) Fingerprint(url string) (map[string]string, error) {
resp, err := w.Client.Get(url)
if err != nil {
return nil, err
}
defer resp.Body.Close()

result := make(map[string]string)
result["server"] = resp.Header.Get("Server")
result["content_type"] = resp.Header.Get("Content-Type")
result["status"] = resp.Status

// Detect CMS
body := make([]byte, 4096)
resp.Body.Read(body)
bodyStr := string(body)

if strings.Contains(bodyStr, "wp-content") || strings.Contains(bodyStr, "wp-includes") {
result["cms"] = "WordPress"
} else if strings.Contains(bodyStr, "Drupal") {
result["cms"] = "Drupal"
} else if strings.Contains(bodyStr, "Joomla") {
result["cms"] = "Joomla"
}

return result, nil
}

func (w *WebOSINT) DetectWAF(url string) (string, error) {
resp, err := w.Client.Get(url + "/?id=1'")
if err != nil {
return "", err
}
defer resp.Body.Close()

if resp.Header.Get("X-Sucuri-ID") != "" {
return "Sucuri", nil
}
if resp.Header.Get("X-Cloudflare") != "" {
return "Cloudflare", nil
}
return "Unknown", nil
}
