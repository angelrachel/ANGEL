package aws

import (
"io"
"net/http"
"time"
)

type AWSCred struct {
Client *http.Client
}

func NewAWSCred() *AWSCred {
return &AWSCred{
Client: &http.Client{Timeout: 5 * time.Second},
}
}

func (a *AWSCred) GetInstanceID() (string, error) {
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

func (a *AWSCred) GetIAMRole() (string, error) {
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

func (a *AWSCred) GetCredentials(role string) (string, error) {
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

func (a *AWSCred) GetUserData() (string, error) {
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

func (a *AWSCred) DumpAllCredentials() (string, error) {
var result string
id, _ := a.GetInstanceID()
result += "Instance ID: " + id + "\n"
role, _ := a.GetIAMRole()
result += "IAM Role: " + role + "\n"
creds, _ := a.GetCredentials(role)
result += "Credentials: " + creds + "\n"
return result, nil
}
