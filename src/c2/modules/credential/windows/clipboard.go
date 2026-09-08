package credential

import "os/exec"

type ClipboardResult struct {
Content string
}

func GetClipboard() ClipboardResult {
cmd := exec.Command("powershell", "-Command", "Get-Clipboard")
out, _ := cmd.Output()
return ClipboardResult{Content: string(out)}
}

func SetClipboard(content string) ClipboardResult {
cmd := exec.Command("powershell", "-Command", "Set-Clipboard -Value '"+content+"'")
cmd.Run()
return ClipboardResult{Content: content}
}
