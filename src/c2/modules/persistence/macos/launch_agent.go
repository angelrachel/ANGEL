package persistence

import (
"os"
"os/exec"
)

type LaunchAgentResult struct {
Label  string
Status string
}

func CreateLaunchAgent(label, command string) LaunchAgentResult {
content := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n<key>Label</key>\n<string>" + label + "</string>\n<key>ProgramArguments</key>\n<array>\n<string>" + command + "</string>\n</array>\n<key>RunAtLoad</key>\n<true/>\n</dict>\n</plist>"
os.WriteFile("/Library/LaunchAgents/"+label+".plist", []byte(content), 0644)
cmd := exec.Command("launchctl", "load", "/Library/LaunchAgents/"+label+".plist")
cmd.Run()
return LaunchAgentResult{Label: label, Status: "success"}
}

func DeleteLaunchAgent(label string) LaunchAgentResult {
cmd := exec.Command("launchctl", "unload", "/Library/LaunchAgents/"+label+".plist")
cmd.Run()
return LaunchAgentResult{Label: label, Status: "success"}
}
