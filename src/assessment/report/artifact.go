package report

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
)

type Artifact struct {
	Format      string `json:"format"`
	ContentType string `json:"content_type"`
	Body        []byte `json:"body"`
	SHA256      string `json:"sha256"`
}

func RenderArtifact(report Report, format string) (Artifact, error) {
	if err := report.Validate(); err != nil {
		return Artifact{}, err
	}
	format = strings.ToLower(strings.TrimSpace(format))
	var body []byte
	contentType := ""
	switch format {
	case "json":
		var err error
		body, err = json.MarshalIndent(report, "", "  ")
		if err != nil {
			return Artifact{}, err
		}
		contentType = "application/json"
	case "markdown", "md":
		body = []byte(report.Markdown())
		format = "markdown"
		contentType = "text/markdown; charset=utf-8"
	case "html":
		rendered, err := renderHTML(report)
		if err != nil {
			return Artifact{}, err
		}
		body = []byte(rendered)
		contentType = "text/html; charset=utf-8"
	default:
		return Artifact{}, fmt.Errorf("unsupported report format: %s", format)
	}
	digest := sha256.Sum256(body)
	return Artifact{Format: format, ContentType: contentType, Body: body, SHA256: hex.EncodeToString(digest[:])}, nil
}

var reportHTML = template.Must(template.New("angel-report").Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><title>{{.Title}}</title></head>
<body><main><h1>{{.Title}}</h1>{{range .Findings}}<article><h2>{{.Title}}</h2><p><strong>Severity:</strong> {{.Severity}}</p><p>{{.Description}}</p><p><strong>Evidence:</strong> {{.Evidence}}</p></article>{{end}}</main></body></html>`))

func renderHTML(report Report) (string, error) {
	var output strings.Builder
	if err := reportHTML.Execute(&output, report); err != nil {
		return "", err
	}
	return output.String(), nil
}
