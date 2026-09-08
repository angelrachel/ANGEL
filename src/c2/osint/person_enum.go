package osint

import (
"net/http"
"strings"
)

type PersonResult struct {
Name   string
Email  string
Source string
}

type PersonScanner struct {
Client http.Client
}

func NewPersonScanner() *PersonScanner {
return &PersonScanner{Client: http.Client{}}
}

func (p *PersonScanner) SearchGoogle(name string) []PersonResult {
var results []PersonResult
url := "https://www.google.com/search?q=" + strings.ReplaceAll(name, " ", "+")
resp, err := p.Client.Get(url)
if err != nil {
return results
}
defer resp.Body.Close()
body := make([]byte, 4096)
resp.Body.Read(body)
results = append(results, PersonResult{Name: name, Email: "", Source: "Google"})
return results
}

func (p *PersonScanner) SearchLinkedIn(name string) []PersonResult {
var results []PersonResult
url := "https://www.google.com/search?q=site:linkedin.com+" + strings.ReplaceAll(name, " ", "+")
resp, err := p.Client.Get(url)
if err != nil {
return results
}
defer resp.Body.Close()
body := make([]byte, 4096)
resp.Body.Read(body)
results = append(results, PersonResult{Name: name, Email: "", Source: "LinkedIn"})
return results
}

func (p *PersonScanner) SearchEmail(email string) []PersonResult {
var results []PersonResult
url := "https://www.google.com/search?q=" + strings.ReplaceAll(email, "@", "%40")
resp, err := p.Client.Get(url)
if err != nil {
return results
}
defer resp.Body.Close()
body := make([]byte, 4096)
resp.Body.Read(body)
results = append(results, PersonResult{Name: "", Email: email, Source: "Google"})
return results
}
