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

func (a *AWSMetadata) GetMetadata(path string) (string, error) {
resp, err := a.Client.Get("http://169.254.169.254/latest/meta-data/" + path)
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

func (a *AWSMetadata) GetInstanceID() (string, error) {
return a.GetMetadata("instance-id")
}

func (a *AWSMetadata) GetInstanceType() (string, error) {
return a.GetMetadata("instance-type")
}

func (a *AWSMetadata) GetHostname() (string, error) {
return a.GetMetadata("hostname")
}

func (a *AWSMetadata) GetMacAddress() (string, error) {
return a.GetMetadata("mac")
}

func (a *AWSMetadata) GetAmiID() (string, error) {
return a.GetMetadata("ami-id")
}
