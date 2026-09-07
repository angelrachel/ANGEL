package network_evasion

import (
"crypto/rand"
"encoding/hex"
"fmt"
"net/http"
"net/url"
"time"
)

type IPRotationEngine struct {
ProxyPool    []string
CurrentProxy int
UserAgents   []string
}

func NewIPRotationEngine() *IPRotationEngine {
return &IPRotationEngine{
ProxyPool: []string{},
UserAgents: []string{
"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Safari/605.1.15",
"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
},
}
}

func (i *IPRotationEngine) AddProxy(proxy string) {
i.ProxyPool = append(i.ProxyPool, proxy)
}

func (i *IPRotationEngine) NextProxy() string {
if len(i.ProxyPool) == 0 {
return ""
}
proxy := i.ProxyPool[i.CurrentProxy]
i.CurrentProxy = (i.CurrentProxy + 1) % len(i.ProxyPool)
return proxy
}

func (i *IPRotationEngine) NextUserAgent() string {
return i.UserAgents[time.Now().UnixNano()%int64(len(i.UserAgents))]
}

func (i *IPRotationEngine) GenerateJitter(base time.Duration, variance time.Duration) time.Duration {
n, _ := rand.Int(rand.Reader, big.NewInt(int64(variance)))
return base + time.Duration(n.Int64())
}

func (i *IPRotationEngine) CreateClient() *http.Client {
proxyStr := i.NextProxy()
transport := &http.Transport{}
if proxyStr != "" {
proxyURL, _ := url.Parse(proxyStr)
transport.Proxy = http.ProxyURL(proxyURL)
}
return &http.Client{
Timeout:   15 * time.Second,
Transport: transport,
}
}

func (i *IPRotationEngine) ExecuteRequest(req *http.Request) (*http.Response, error) {
client := i.CreateClient()
req.Header.Set("User-Agent", i.NextUserAgent())
time.Sleep(i.GenerateJitter(500*time.Millisecond, 300*time.Millisecond))
return client.Do(req)
}

func (i *IPRotationEngine) RotateWithDelay(delay time.Duration) {
time.Sleep(delay)
}

func (i *IPRotationEngine) GetProxyList() []string {
return i.ProxyPool
}

func (i *IPRotationEngine) GenerateRandomIPv4() string {
b := make([]byte, 4)
rand.Read(b)
return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}

func (i *IPRotationEngine) GenerateRandomIPv6() string {
b := make([]byte, 16)
rand.Read(b)
return hex.EncodeToString(b)
}
