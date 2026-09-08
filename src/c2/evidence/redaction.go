package evidence

import (
"regexp"
)

type Redactor struct {
Patterns []*regexp.Regexp
}

func NewRedactor() *Redactor {
r := &Redactor{}
r.Patterns = append(r.Patterns, regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[:=]\s*\S+`))
r.Patterns = append(r.Patterns, regexp.MustCompile(`(?i)(api[_-]?key|secret|token)\s*[:=]\s*\S+`))
r.Patterns = append(r.Patterns, regexp.MustCompile(`(?i)authorization:\s*bearer\s+\S+`))
return r
}

func (r *Redactor) Redact(data string) string {
for _, pattern := range r.Patterns {
data = pattern.ReplaceAllString(data, "[REDACTED]")
}
return data
}

func (r *Redactor) RedactPII(data string) string {
email := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
phone := regexp.MustCompile(`\+?[0-9]{10,15}`)
data = email.ReplaceAllString(data, "[EMAIL_REDACTED]")
return phone.ReplaceAllString(data, "[PHONE_REDACTED]")
}

func (r *Redactor) RedactAll(data string) string {
data = r.Redact(data)
data = r.RedactPII(data)
return data
}
