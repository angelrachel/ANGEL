package evidence

import (
"encoding/base64"
"io"
"net/http"
"time"
)

type EvidenceCollector struct {
Requests  []RequestRecord
Responses []ResponseRecord
}

type RequestRecord struct {
Timestamp   int64
Method      string
URL         string
Headers     map[string]string
Body        string
}

type ResponseRecord struct {
Timestamp   int64
StatusCode  int
Headers     map[string]string
Body        string
}

func NewEvidenceCollector() *EvidenceCollector {
return &EvidenceCollector{
Requests:  make([]RequestRecord, 0),
Responses: make([]ResponseRecord, 0),
}
}

func (e *EvidenceCollector) CaptureRequest(req *http.Request) error {
body, err := io.ReadAll(req.Body)
if err != nil {
return err
}
record := RequestRecord{
Timestamp: time.Now().Unix(),
Method:    req.Method,
URL:       req.URL.String(),
Headers:   make(map[string]string),
Body:      string(body),
}
for key, values := range req.Header {
if len(values) > 0 {
record.Headers[key] = values[0]
}
}
e.Requests = append(e.Requests, record)
return nil
}

func (e *EvidenceCollector) CaptureResponse(resp *http.Response) error {
body, err := io.ReadAll(resp.Body)
if err != nil {
return err
}
record := ResponseRecord{
Timestamp:  time.Now().Unix(),
StatusCode: resp.StatusCode,
Headers:    make(map[string]string),
Body:       string(body),
}
for key, values := range resp.Header {
if len(values) > 0 {
record.Headers[key] = values[0]
}
}
e.Responses = append(e.Responses, record)
return nil
}

func (e *EvidenceCollector) TakeScreenshot() string {
// In real implementation, would capture screen
return "screenshot_data"
}

func (e *EvidenceCollector) CaptureBeforeAfter(before, after string) DiffRecord {
return DiffRecord{
Before: before,
After:  after,
}
}

type DiffRecord struct {
Before string
After  string
}
