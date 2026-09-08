package osint

import (
"net/http"
"strings"
)

type WebResult struct {
Target  string
Headers []string
Title   string
Status  int
}

type WebScanner struct {
Client http.Client
}

func NewWebScanner() *WebScanner {
return &WebScanner{Client: http.Client{}}
}

func (w *WebScanner) Scan(target string) WebResult {
result := WebResult{Target: target}
resp, err := w.Client.Get(target)
if err != nil {
return result
}
defer resp.Body.Close()
result.Status = resp.StatusCode
for k, v := range resp.Header {
result.Headers = append(result.Headers, k+": "+strings.Join(v, ", "))
}
body := make([]byte, 1024)
resp.Body.Read(body)
content := string(body)
if strings.Contains(content, "<title>") {
start := strings.Index(content, "<title>") + 7
end := strings.Index(content, "</title>")
if end > start {
result.Title = content[start:end]
}
}
return result
}

func (w *WebScanner) FetchHeaders(target string) []string {
resp, err := w.Client.Get(target)
if err != nil {
return []string{}
}
defer resp.Body.Close()
var headers []string
for k, v := range resp.Header {
headers = append(headers, k+": "+strings.Join(v, ", "))
}
return headers
}
