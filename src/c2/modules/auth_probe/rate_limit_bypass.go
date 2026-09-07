package auth_probe

import (
"net/http"
"time"
)

type RateLimitBypass struct {
UserAgents    []string
ProxyList     []string
CurrentAgent  int
CurrentProxy  int
RandomDelay   time.Duration
Client        *http.Client
}

func NewRateLimitBypass() *RateLimitBypass {
return &RateLimitBypass{
UserAgents: []string{
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
},
ProxyList:    []string{},
CurrentAgent: 0,
CurrentProxy: 0,
RandomDelay:  500 * time.Millisecond,
}
}

func (r *RateLimitBypass) NextUserAgent() string {
agent := r.UserAgents[r.CurrentAgent]
r.CurrentAgent = (r.CurrentAgent + 1) % len(r.UserAgents)
return agent
}

func (r *RateLimitBypass) NextProxy() string {
if len(r.ProxyList) == 0 {
return ""
}
proxy := r.ProxyList[r.CurrentProxy]
r.CurrentProxy = (r.CurrentProxy + 1) % len(r.ProxyList)
return proxy
}

func (r *RateLimitBypass) AddProxy(proxy string) {
r.ProxyList = append(r.ProxyList, proxy)
}

func (r *RateLimitBypass) SetRandomDelay(delay time.Duration) {
r.RandomDelay = delay
}

func (r *RateLimitBypass) GenerateRandomDelay() time.Duration {
return r.RandomDelay + time.Duration(time.Now().UnixNano()%int64(r.RandomDelay))
}

func (r *RateLimitBypass) CreateClient() *http.Client {
return &http.Client{Timeout: 15 * time.Second}
}

func (r *RateLimitBypass) GetAllUserAgents() []string {
return r.UserAgents
}

func (r *RateLimitBypass) GetAllProxies() []string {
return r.ProxyList
}
