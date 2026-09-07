package collector

import (
"os"
"time"
)

type WebcamCollector struct{}

func NewWebcamCollector() *WebcamCollector {
return &WebcamCollector{}
}

func (w *WebcamCollector) Capture() (string, error) {
filename := "webcam_" + time.Now().Format("20060102_150405") + ".jpg"
return filename, os.WriteFile(filename, []byte("webcam_data"), 0644)
}

func (w *WebcamCollector) Stream(duration int) (string, error) {
filename := "webcam_stream_" + time.Now().Format("20060102_150405") + ".avi"
return filename, os.WriteFile(filename, []byte("stream_data"), 0644)
}
