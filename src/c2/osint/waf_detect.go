package osint

import (
"io"
"net/http"
"strings"
"time"
)

type WAFDetect struct {
Client *http.Client
}

func NewWAFDetect() *WAFDetect {
return &WAFDetect{
Client: &http.Client{Timeout: 10 * time.Second},
}
}

func (w *WAFDetect) Detect(url string) (string, error) {
payload := "' OR 1=1 -- -"
resp, err := w.Client.Get(url + "?id=" + payload)
if err != nil {
return "", err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil {
return "", err
}
content := string(body)
if strings.Contains(content, "Cloudflare") {
return "Cloudflare", nil
}
if strings.Contains(content, "Sucuri") {
return "Sucuri", nil
}
if strings.Contains(content, "ModSecurity") {
return "ModSecurity", nil
}
if strings.Contains(content, "AWS WAF") {
return "AWS WAF", nil
}
return "No WAF Detected", nil
}

func (w *WAFDetect) DetectWithCookie(url string) (string, error) {
req, _ := http.NewRequest("GET", url, nil)
req.Header.Set("Cookie", "id=1' OR '1'='1")
resp, err := w.Client.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil {
return "", err
}
content := string(body)
if strings.Contains(content, "Cloudflare") {
return "Cloudflare", nil
}
if strings.Contains(content, "Sucuri") {
return "Sucuri", nil
}
return "No WAF Detected", nil
}

func (w *WAFDetect) GetWAFType(url string) (string, error) {
return w.Detect(url)
}
