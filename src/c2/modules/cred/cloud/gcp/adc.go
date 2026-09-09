package gcp

import (
	"os"
	"strings"
)

type ADC struct {
	Type         string
	ProjectID    string
	ClientEmail  string
	PrivateKeyID string
}

func GetADCFromFile(path string) ADC {
	data, err := os.ReadFile(path)
	if err != nil {
		return ADC{}
	}
	var adc ADC
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "\"") {
			parts := strings.SplitN(line, "\"", 2)
			if strings.HasPrefix(parts[0], "type") {
				adc.Type = strings.TrimSuffix(parts[1], "\"")
			} else if strings.HasPrefix(parts[0], "project_id") {
				adc.ProjectID = strings.TrimSuffix(parts[1], "\"")
			} else if strings.HasPrefix(parts[0], "client_email") {
				adc.ClientEmail = strings.TrimSuffix(parts[1], "\"")
			}
		}
	}
	return adc
}
