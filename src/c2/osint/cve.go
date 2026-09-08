package osint

import (
"io"
"net/http"
"strings"
"time"
)

type CVEDetect struct {
Client *http.Client
}

func NewCVEDetect() *CVEDetect {
return &CVEDetect{
Client: &http.Client{Timeout: 10 * time.Second},
}
}

func (c *CVEDetect) Detect(url string) string {
resp, err := c.Client.Get(url)
if err != nil {
return ""
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
content := string(body)
if strings.Contains(content, "CVE-2021-44228") {
return "Log4j (CVE-2021-44228)"
}
if strings.Contains(content, "CVE-2018-7600") {
return "Drupalgeddon2 (CVE-2018-7600)"
}
if strings.Contains(content, "CVE-2019-0708") {
return "BlueKeep (CVE-2019-0708)"
}
return "No known CVE Detected"
}

func (c *CVEDetect) DetectWithPayload(url string) string {
resp, err := c.Client.Get(url + "/?=${jndi:ldap://127.0.0.1/a}")
if err != nil {
return ""
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
content := string(body)
if strings.Contains(content, "CVE-2021-44228") {
return "Log4j (CVE-2021-44228) - Confirmed"
}
return "No CVE Detected"
}
