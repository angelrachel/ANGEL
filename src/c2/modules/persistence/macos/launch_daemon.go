package persistence

import "os/exec"

type LaunchDaemonResult struct {
Label  string
Status string
}

func CreateLaunchDaemon(label, command string) LaunchDaemonResult {
cmd := exec.Command("bash", "-c", "echo '<?xml version=\"1.0\" encoding=\"UTF-8\"?>' > /Library/LaunchDaemons/"+label+".plist && echo '<plist version=\"1.0\"><dict><key>Label</key><string>"+label+"</string><key>ProgramArguments</key><array><string>"+command+"</string></array><key>RunAtLoad</key><true/></dict></plist>' >> /Library/LaunchDaemons/"+label+".plist")
cmd.Run()
return LaunchDaemonResult{Label: label, Status: "success"}
}

func DeleteLaunchDaemon(label string) LaunchDaemonResult {
cmd := exec.Command("launchctl", "unload", "/Library/LaunchDaemons/"+label+".plist")
cmd.Run()
return LaunchDaemonResult{Label: label, Status: "success"}
}
