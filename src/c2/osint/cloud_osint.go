package osint

import (
"io"
"net/http"
"time"
)

type CloudOSINT struct {
Client *http.Client
}

func NewCloudOSINT() *CloudOSINT {
return &CloudOSINT{
Client: &http.Client{Timeout: 5 * time.Second},
}
}

func (c *CloudOSINT) GetAWSMetadata() (string, error) {
resp, err := c.Client.Get("http://169.254.169.254/latest/meta-data/")
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

func (c *CloudOSINT) GetGCPMetadata() (string, error) {
req, _ := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/", nil)
req.Header.Set("Metadata-Flavor", "Google")
resp, err := c.Client.Do(req)
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

func (c *CloudOSINT) GetAzureMetadata() (string, error) {
req, _ := http.NewRequest("GET", "http://169.254.169.254/metadata/instance?api-version=2021-02-01", nil)
req.Header.Set("Metadata", "true")
resp, err := c.Client.Do(req)
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

func (c *CloudOSINT) DumpAll() (string, error) {
var result string
aws, _ := c.GetAWSMetadata()
result += "AWS: " + aws + "\n"
gcp, _ := c.GetGCPMetadata()
result += "GCP: " + gcp + "\n"
azure, _ := c.GetAzureMetadata()
result += "Azure: " + azure + "\n"
return result, nil
}
