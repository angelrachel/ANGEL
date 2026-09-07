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

func (a *AWSMetadata) GetInstanceID() (string, error) {
resp, err := a.Client.Get("http://169.254.169.254/latest/meta-data/instance-id")
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

func (a *AWSMetadata) GetIAMRole() (string, error) {
resp, err := a.Client.Get("http://169.254.169.254/latest/meta-data/iam/security-credentials/")
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

func (a *AWSMetadata) GetCredentials(role string) (string, error) {
resp, err := a.Client.Get("http://169.254.169.254/latest/meta-data/iam/security-credentials/" + role)
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

func (a *AWSMetadata) GetUserData() (string, error) {
resp, err := a.Client.Get("http://169.254.169.254/latest/user-data")
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
