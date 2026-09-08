package osint

import (
"net/http"
"strings"
)

type WAFResult struct {
Target string
WAF    string
}

type WAFDetector struct {
Client http.Client
}

func NewWAFDetector() *WAFDetector {
return &WAFDetector{Client: http.Client{}}
}

func (w *WAFDetector) Detect(target string) WAFResult {
result := WAFResult{Target: target, WAF: "None"}
payload := "?id=1' OR '1'='1"
resp, err := w.Client.Get(target + payload)
if err != nil {
return result
}
defer resp.Body.Close()
body := make([]byte, 512)
resp.Body.Read(body)

wafs := []string{"Cloudflare", "Sucuri", "ModSecurity", "AWS WAF", "Imperva", "F5 BIG-IP"}
for _, waf := range wafs {
if strings.Contains(string(body), waf) || strings.Contains(strings.ToLower(string(resp.Header.Get("Server"))), strings.ToLower(waf)) {
result.WAF = waf
return result
}
}
return result
}
