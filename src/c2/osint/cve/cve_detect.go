package cve

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

func (c *CVEDetect) Detect(url string) (string, error) {
resp, err := c.Client.Get(url)
if err != nil {
return "", err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil {
return "", err
}
content := string(body)
if strings.Contains(content, "CVE-2021-44228") {
return "Log4j (CVE-2021-44228)", nil
}
if strings.Contains(content, "CVE-2018-7600") {
return "Drupalgeddon2 (CVE-2018-7600)", nil
}
if strings.Contains(content, "CVE-2019-0708") {
return "BlueKeep (CVE-2019-0708)", nil
}
return "No known CVE Detected", nil
}

func (c *CVEDetect) DetectWithPayload(url string) (string, error) {
resp, err := c.Client.Get(url + "/?=${jndi:ldap://127.0.0.1/a}")
if err != nil {
return "", err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil {
return "", err
}
content := string(body)
if strings.Contains(content, "CVE-2021-44228") {
return "Log4j (CVE-2021-44228) - Confirmed", nil
}
return "No CVE Detected", nil
}

func (c *CVEDetect) DetectAll(url string) (string, error) {
var result string
detect, _ := c.Detect(url)
result += detect + "\n"
payload, _ := c.DetectWithPayload(url)
result += payload + "\n"
return result, nil
}
