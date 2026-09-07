package evidence

import (
"regexp"
"strings"
)

type Redaction struct {
Patterns []string
}

func NewRedaction() *Redaction {
return &Redaction{
Patterns: []string{
`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}`,
`\d{3}-\d{2}-\d{4}`,
`[A-Z]{2,5}-[A-Z0-9]{4,20}`,
`-----BEGIN [A-Z]+ PRIVATE KEY-----`,
`[0-9a-f]{32,64}`,
`Bearer\s+[A-Za-z0-9._-]+`,
`Authorization:\s*[A-Za-z0-9._-]+`,
},
}
}

func (r *Redaction) Redact(data string) string {
result := data
for _, pattern := range r.Patterns {
re := regexp.MustCompile(pattern)
result = re.ReplaceAllString(result, "[REDACTED]")
}
return result
}

func (r *Redaction) RedactCustom(data string, pattern string) string {
re := regexp.MustCompile(pattern)
return re.ReplaceAllString(data, "[REDACTED]")
}

func (r *Redaction) IsRedacted(data string) bool {
return strings.Contains(data, "[REDACTED]")
}

func (r *Redaction) AddPattern(pattern string) {
r.Patterns = append(r.Patterns, pattern)
}
