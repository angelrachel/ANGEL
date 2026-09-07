package network_evasion

import (
"crypto/tls"
"net/http"
"time"
)

type TrafficMorphing struct {
Client *http.Client
}

func NewTrafficMorphing() *TrafficMorphing {
return &TrafficMorphing{
Client: &http.Client{
Timeout: 15 * time.Second,
Transport: &http.Transport{
TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
},
},
}
}

func (t *TrafficMorphing) MorphRequest(req *http.Request) *http.Request {
req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
req.Header.Set("Accept-Language", "en-US,en;q=0.5")
req.Header.Set("Accept-Encoding", "gzip, deflate, br")
req.Header.Set("Connection", "keep-alive")
req.Header.Set("Upgrade-Insecure-Requests", "1")
return req
}

func (t *TrafficMorphing) ExecuteMorphed(url string, method string) (*http.Response, error) {
var req *http.Request
var err error
if method == "GET" {
req, err = http.NewRequest("GET", url, nil)
} else {
req, err = http.NewRequest("POST", url, nil)
}
if err != nil {
return nil, err
}
t.MorphRequest(req)
return t.Client.Do(req)
}

func (t *TrafficMorphing) SimulateBrowserTraffic(url string) (*http.Response, error) {
req, _ := http.NewRequest("GET", url, nil)
req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
return t.Client.Do(req)
}

func (t *TrafficMorphing) GenerateSleepJitter() time.Duration {
return t.GenerateJitter(200*time.Millisecond, 500*time.Millisecond)
}

func (t *TrafficMorphing) GetAllHeaders() http.Header {
return http.Header{
"Accept":          {"text/html"},
"Accept-Language": {"en-US"},
}
}
