package evidence

import "regexp"

type Redactor struct {
	Regexes []*regexp.Regexp
}

func NewRedactor() *Redactor {
	patterns := []string{
		`(?i)password\s*=\s*["'][^"']*["']`,
		`(?i)api[_-]?key\s*=\s*["'][^"']*["']`,
		`(?i)secret\s*=\s*["'][^"']*["']`,
		`Bearer\s+[A-Za-z0-9._~+/=-]+`,
		`\b\d{3}-\d{2}-\d{4}\b`,
	}
	var regexes []*regexp.Regexp
	for _, pattern := range patterns {
		re, _ := regexp.Compile(pattern)
		regexes = append(regexes, re)
	}
	return &Redactor{Regexes: regexes}
}

func (r *Redactor) Redact(input string) string {
	if r == nil {
		return input
	}
	for _, re := range r.Regexes {
		if re == nil {
			continue
		}
		input = re.ReplaceAllString(input, "[REDACTED]")
	}
	return input
}
