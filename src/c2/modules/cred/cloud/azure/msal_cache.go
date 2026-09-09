package azure

import (
	"os"
	"strings"
)

type MSALCache struct {
	ClientID string
	Account  string
	Token    string
}

func GetMSALCache(path string) MSALCache {
	data, err := os.ReadFile(path)
	if err != nil {
		return MSALCache{}
	}
	var cache MSALCache
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key == "client_id" {
				cache.ClientID = value
			} else if key == "account" {
				cache.Account = value
			} else if key == "token" {
				cache.Token = value
			}
		}
	}
	return cache
}

func GetTokenFromEnv() string {
	return os.Getenv("AZURE_ACCESS_TOKEN")
}
