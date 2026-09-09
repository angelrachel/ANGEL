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

func (w *WAFDetect) Detect(url string) string {
	resp, err := w.Client.Get(url + "/?id=1' OR '1'='1")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	content := string(body)
	if strings.Contains(content, "Cloudflare") {
		return "Cloudflare"
	}
	if strings.Contains(content, "Sucuri") {
		return "Sucuri"
	}
	if strings.Contains(content, "ModSecurity") {
		return "ModSecurity"
	}
	return "No WAF Detected"
}

func (w *WAFDetect) DetectWithCookie(url string) string {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", "id=1' OR '1'='1")
	resp, err := w.Client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	content := string(body)
	if strings.Contains(content, "Cloudflare") {
		return "Cloudflare"
	}
	return "No WAF Detected"
}
