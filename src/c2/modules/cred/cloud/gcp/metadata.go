package gcp

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

func (g *GCPMetadata) GetProjectID() string {
	req, _ := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := g.Client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func (g *GCPMetadata) GetInstanceID() string {
	req, _ := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/instance/id", nil)
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := g.Client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func (g *GCPMetadata) GetToken() string {
	req, _ := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token", nil)
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := g.Client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}
