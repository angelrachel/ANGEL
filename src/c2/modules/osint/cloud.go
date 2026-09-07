package osint

import (
"encoding/xml"
"fmt"
"io"
"net/http"
"strings"
)

type CloudOSINT struct{}

func NewCloudOSINT() *CloudOSINT {
return &CloudOSINT{}
}

func (c *CloudOSINT) EnumerateAWSBucket(bucket string) ([]string, error) {
url := "https://" + bucket + ".s3.amazonaws.com/"
resp, err := http.Get(url)
if err != nil {
return nil, err
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
return nil, err
}

var result struct {
Contents []struct {
Key string `xml:"Key"`
} `xml:"Contents"`
}
xml.Unmarshal(body, &result)

var files []string
for _, item := range result.Contents {
files = append(files, item.Key)
}
return files, nil
}

func (c *CloudOSINT) EnumerateAzureBlob(account string) ([]string, error) {
url := "https://" + account + ".blob.core.windows.net/"
resp, err := http.Get(url)
if err != nil {
return nil, err
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
return nil, err
}

// Parse XML response for Azure Blob
var containers []string
lines := strings.Split(string(body), "\n")
for _, line := range lines {
if strings.Contains(line, "<Name>") {
// Extract container name
continue
}
}
return containers, nil
}
