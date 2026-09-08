package aws

import (
"os"
"strings"
)

type AWSConfig struct {
AccessKey    string
SecretKey    string
SessionToken string
Region       string
}

func GetCredentialsFromEnv() AWSConfig {
return AWSConfig{
AccessKey:    os.Getenv("AWS_ACCESS_KEY_ID"),
SecretKey:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
SessionToken: os.Getenv("AWS_SESSION_TOKEN"),
Region:       os.Getenv("AWS_REGION"),
}
}

func GetCredentialsFromFile(path string) AWSConfig {
data, err := os.ReadFile(path)
if err != nil {
return AWSConfig{}
}
var config AWSConfig
for _, line := range strings.Split(string(data), "\n") {
line = strings.TrimSpace(line)
if strings.Contains(line, "=") {
parts := strings.SplitN(line, "=", 2)
key := strings.TrimSpace(parts[0])
value := strings.TrimSpace(parts[1])
if key == "aws_access_key_id" {
config.AccessKey = value
} else if key == "aws_secret_access_key" {
config.SecretKey = value
} else if key == "aws_session_token" {
config.SessionToken = value
}
}
}
return config
}
