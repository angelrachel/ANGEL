package azure

import (
"io"
"net/http"
"time"
)

type AzureMetadata struct {
Client *http.Client
}

func NewAzureMetadata() *AzureMetadata {
return &AzureMetadata{
Client: &http.Client{Timeout: 5 * time.Second},
}
}

func (a *AzureMetadata) GetSubscriptionID() string {
req, _ := http.NewRequest("GET", "http://169.254.169.254/metadata/instance/compute/subscriptionId?api-version=2021-02-01&format=text", nil)
req.Header.Set("Metadata", "true")
resp, err := a.Client.Do(req)
if err != nil {
return ""
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
return string(body)
}

func (a *AzureMetadata) GetLocation() string {
req, _ := http.NewRequest("GET", "http://169.254.169.254/metadata/instance/compute/location?api-version=2021-02-01&format=text", nil)
req.Header.Set("Metadata", "true")
resp, err := a.Client.Do(req)
if err != nil {
return ""
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
return string(body)
}

func (a *AzureMetadata) GetVMName() string {
req, _ := http.NewRequest("GET", "http://169.254.169.254/metadata/instance/compute/name?api-version=2021-02-01&format=text", nil)
req.Header.Set("Metadata", "true")
resp, err := a.Client.Do(req)
if err != nil {
return ""
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
return string(body)
}
