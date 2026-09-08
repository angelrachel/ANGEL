package osint

import (
"net/http"
"strings"
)

type CloudResult struct {
Target string
Data   string
}

type CloudScanner struct {
Client http.Client
}

func NewCloudScanner() *CloudScanner {
return &CloudScanner{Client: http.Client{}}
}

func (c *CloudScanner) AWSMetadata() CloudResult {
result := CloudResult{Target: "169.254.169.254"}
resp, err := c.Client.Get("http://169.254.169.254/latest/meta-data/")
if err != nil {
return result
}
defer resp.Body.Close()
body := make([]byte, 2048)
resp.Body.Read(body)
result.Data = string(body)
return result
}

func (c *CloudScanner) AzureMetadata() CloudResult {
result := CloudResult{Target: "169.254.169.254"}
req, _ := http.NewRequest("GET", "http://169.254.169.254/metadata/instance?api-version=2021-02-01", nil)
req.Header.Set("Metadata", "true")
resp, err := c.Client.Do(req)
if err != nil {
return result
}
defer resp.Body.Close()
body := make([]byte, 2048)
resp.Body.Read(body)
result.Data = string(body)
return result
}

func (c *CloudScanner) BucketEnum(bucketName string) CloudResult {
result := CloudResult{Target: bucketName}
resp, err := c.Client.Get("https://" + bucketName + ".s3.amazonaws.com")
if err != nil {
return result
}
defer resp.Body.Close()
body := make([]byte, 2048)
resp.Body.Read(body)
if strings.Contains(string(body), "ListBucket") {
result.Data = string(body)
}
return result
}
