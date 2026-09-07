package cloud

import (
"io"
"net/http"
"time"
)

type GCPMetadata struct {
Client *http.Client
}

func NewGCPMetadata() *GCPMetadata {
return &GCPMetadata{
Client: &http.Client{Timeout: 5 * time.Second},
}
}

func (g *GCPMetadata) GetProjectID() (string, error) {
req, _ := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)
req.Header.Set("Metadata-Flavor", "Google")
resp, err := g.Client.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil {
return "", err
}
return string(body), nil
}

func (g *GCPMetadata) GetInstanceID() (string, error) {
req, _ := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/instance/id", nil)
req.Header.Set("Metadata-Flavor", "Google")
resp, err := g.Client.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil {
return "", err
}
return string(body), nil
}

func (g *GCPMetadata) GetServiceAccount() (string, error) {
req, _ := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/email", nil)
req.Header.Set("Metadata-Flavor", "Google")
resp, err := g.Client.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil {
return "", err
}
return string(body), nil
}

func (g *GCPMetadata) GetToken() (string, error) {
req, _ := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token", nil)
req.Header.Set("Metadata-Flavor", "Google")
resp, err := g.Client.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil {
return "", err
}
return string(body), nil
}
