package redaction

import (
"regexp"
)

type SecretFilter struct {
Patterns []*regexp.Regexp
}

func NewSecretFilter() *SecretFilter {
f := &SecretFilter{}
f.Patterns = append(f.Patterns, regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[:=]\s*\S+`))
f.Patterns = append(f.Patterns, regexp.MustCompile(`(?i)(api[_-]?key|secret|token)\s*[:=]\s*\S+`))
f.Patterns = append(f.Patterns, regexp.MustCompile(`(?i)authorization:\s*bearer\s+\S+`))
return f
}

func (s *SecretFilter) Filter(data string) string {
for _, p := range s.Patterns {
data = p.ReplaceAllString(data, "[REDACTED]")
}
return data
}

func (s *SecretFilter) FilterPII(data string) string {
email := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
phone := regexp.MustCompile(`\+?[0-9]{10,15}`)
return email.ReplaceAllString(phone.ReplaceAllString(data, "[PHONE_REDACTED]"), "[EMAIL_REDACTED]")
}

func (s *SecretFilter) FilterAll(data string) string {
data = s.Filter(data)
data = s.FilterPII(data)
return data
}
