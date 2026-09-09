package credential

import (
	"bufio"
	"os"
	"strings"
)

type CloudCredential struct {
	Provider string
	Key      string
	Value    string
}

func ReadAWS(path string) []CloudCredential {
	var creds []CloudCredential
	file, err := os.Open(path)
	if err != nil {
		return creds
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			creds = append(creds, CloudCredential{Provider: "aws", Key: parts[0], Value: parts[1]})
		}
	}
	return creds
}

func ReadAzure(path string) []CloudCredential {
	var creds []CloudCredential
	file, err := os.Open(path)
	if err != nil {
		return creds
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			creds = append(creds, CloudCredential{Provider: "azure", Key: parts[0], Value: parts[1]})
		}
	}
	return creds
}

func ReadGCP(path string) []CloudCredential {
	var creds []CloudCredential
	file, err := os.Open(path)
	if err != nil {
		return creds
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			creds = append(creds, CloudCredential{Provider: "gcp", Key: parts[0], Value: parts[1]})
		}
	}
	return creds
}
