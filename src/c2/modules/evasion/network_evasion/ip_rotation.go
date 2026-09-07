//go:build windows

package network_evasion

import (
"crypto/rand"
"fmt"
"math/big"
"net/http"
"net/url"
"time"
)

// GenerateRandomIP menghasilkan IP acak baru.
func GenerateRandomIP() string {
b := make([]byte, 4)
rand.Read(b)
return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}

// GenerateJitter menghasilkan delay acak antara minimum dan maksimum.
func GenerateJitter(min, max time.Duration) time.Duration {
diff := max - min
if diff <= 0 {
return min
}
n, _ := rand.Int(rand.Reader, big.NewInt(int64(diff)))
return min + time.Duration(n.Int64())
}

// RotateWithDelay melakukan delay yang menyesuaikan durasi 1-3 detik.
func RotateWithDelay() {
// Base delay 1 detik dengan jitter maksimal 2 detik (total 1-3 detik)
delay := GenerateJitter(1*time.Second, 3*time.Second)
time.Sleep(delay)
}

// CreateProxyClient membuat client dengan proxy baru (auto-rotate IP).
func CreateProxyClient(proxy string) (*http.Client, error) {
proxyURL, err := url.Parse(proxy)
if err != nil {
return nil, err
}
return &http.Client{
Timeout:   15 * time.Second,
Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
}, nil
}

// AutoRotateProxy melakukan rotate proxy dengan delay 1-3 detik.
func AutoRotateProxy(proxies []string) (*http.Client, error) {
randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(proxies))))
client, err := CreateProxyClient(proxies[randomIndex.Int64()])
if err != nil {
return nil, err
}
RotateWithDelay()
return client, nil
}
