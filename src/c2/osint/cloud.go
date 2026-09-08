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

func (c *CloudOSINT) GetAWSMetadata() string {
resp, err := c.Client.Get("http://169.254.169.254/latest/meta-data/")
if err != nil {
return ""
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
return string(body)
}

func (c *CloudOSINT) GetGCPMetadata() string {
req, _ := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/", nil)
req.Header.Set("Metadata-Flavor", "Google")
resp, err := c.Client.Do(req)
if err != nil {
return ""
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
return string(body)
}

func (c *CloudOSINT) GetAzureMetadata() string {
req, _ := http.NewRequest("GET", "http://169.254.169.254/metadata/instance?api-version=2021-02-01", nil)
req.Header.Set("Metadata", "true")
resp, err := c.Client.Do(req)
if err != nil {
return ""
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
return string(body)
}
