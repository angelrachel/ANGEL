package aws

import (
	"io"
	"net/http"
	"time"
)

type AWSMetadata struct {
	Client *http.Client
}

func NewAWSMetadata() *AWSMetadata {
	return &AWSMetadata{
		Client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (a *AWSMetadata) GetInstanceID() string {
	resp, err := a.Client.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func (a *AWSMetadata) GetIAMRole() string {
	resp, err := a.Client.Get("http://169.254.169.254/latest/meta-data/iam/security-credentials/")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func (a *AWSMetadata) GetCredentials(role string) string {
	resp, err := a.Client.Get("http://169.254.169.254/latest/meta-data/iam/security-credentials/" + role)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}
