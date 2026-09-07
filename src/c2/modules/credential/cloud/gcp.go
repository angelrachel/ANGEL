package cloud

import (
"os"
)

type GCP struct{}

func NewGCP() *GCP {
return &GCP{}
}

func (g *GCP) ExtractADC() (string, error) {
path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
if path == "" {
path = os.Getenv("HOME") + "/.config/gcloud/application_default_credentials.json"
}
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}

func (g *GCP) ExtractGCloudConfig() (string, error) {
path := os.Getenv("HOME") + "/.config/gcloud/configurations/config_default"
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}
