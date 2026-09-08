package osint

import (
"net/http"
"strings"
)

type TechResult struct {
Target    string
CMS       string
Server    string
Framework string
Headers   []string
}

type TechScanner struct {
Client http.Client
}

func NewTechScanner() *TechScanner {
return &TechScanner{Client: http.Client{}}
}

func (t *TechScanner) Scan(target string) TechResult {
result := TechResult{Target: target}
resp, err := t.Client.Get(target)
if err != nil {
return result
}
defer resp.Body.Close()

body := make([]byte, 2048)
resp.Body.Read(body)
bodyStr := string(body)

cmsList := []string{"WordPress", "Joomla", "Drupal", "Magento", "Shopify"}
serverList := []string{"nginx", "Apache", "IIS", "Caddy", "LiteSpeed"}
frameworkList := []string{"Laravel", "React", "Vue.js", "Angular", "Next.js", "Nuxt.js"}

for _, cms := range cmsList {
if strings.Contains(strings.ToLower(bodyStr), strings.ToLower(cms)) {
result.CMS = cms
break
}
}
for _, server := range serverList {
if strings.Contains(strings.ToLower(resp.Header.Get("Server")), strings.ToLower(server)) {
result.Server = server
break
}
}
for _, framework := range frameworkList {
if strings.Contains(strings.ToLower(bodyStr), strings.ToLower(framework)) {
result.Framework = framework
break
}
}

for k, v := range resp.Header {
result.Headers = append(result.Headers, k+": "+strings.Join(v, ", "))
}

return result
}
