package report_gen

import (
"fmt"
"os"
"strings"
"time"
)

type PDFReport struct {
Title       string
Description string
Findings    []Finding
GeneratedAt time.Time
}

type Finding struct {
ID          string
Severity    string
Title       string
Description string
Remediation string
}

func NewPDFReport(title, description string) *PDFReport {
return &PDFReport{
Title:       title,
Description: description,
GeneratedAt: time.Now(),
}
}

func (p *PDFReport) AddFinding(f Finding) {
p.Findings = append(p.Findings, f)
}

func (p *PDFReport) GenerateMinimalPDF() string {
var sb strings.Builder
sb.WriteString("%PDF-1.4\n")
sb.WriteString("1 0 obj << /Type /Catalog /Pages 2 0 R >> endobj\n")
sb.WriteString("2 0 obj << /Type /Pages /Kids [3 0 R] /Count 1 >> endobj\n")
sb.WriteString("3 0 obj << /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >> endobj\n")
sb.WriteString("4 0 obj << /Length ")

var content strings.Builder
content.WriteString("BT /F1 24 Tf 50 750 Td (ANGEL RED TEAM REPORT) Tj ET\n")
content.WriteString("BT /F1 12 Tf 50 720 Td (Generated: " + p.GeneratedAt.Format(time.RFC3339) + ") Tj ET\n")
content.WriteString("BT /F1 12 Tf 50 700 Td (" + p.Title + ") Tj ET\n")
content.WriteString("BT /F1 10 Tf 50 680 Td (" + p.Description + ") Tj ET\n")

yPos := 640
for _, f := range p.Findings {
content.WriteString(fmt.Sprintf("BT /F1 12 Tf 50 %d Td (Finding: %s - %s) Tj ET\n", yPos, f.ID, f.Title))
yPos -= 20
content.WriteString(fmt.Sprintf("BT /F1 10 Tf 50 %d Td (Severity: %s) Tj ET\n", yPos, f.Severity))
yPos -= 20
content.WriteString(fmt.Sprintf("BT /F1 10 Tf 50 %d Td (Desc: %s) Tj ET\n", yPos, f.Description))
yPos -= 20
content.WriteString(fmt.Sprintf("BT /F1 10 Tf 50 %d Td (Remediation: %s) Tj ET\n", yPos, f.Remediation))
yPos -= 40
}

contentStr := content.String()
sb.WriteString(fmt.Sprintf("%d", len(contentStr)))
sb.WriteString(" >> stream\n")
sb.WriteString(contentStr)
sb.WriteString("endstream endobj\n")
sb.WriteString("5 0 obj << /Type /Font /Subtype /Type1 /BaseFont /Helvetica >> endobj\n")
sb.WriteString("xref\n")
sb.WriteString("0 6\n")
sb.WriteString("0000000000 65535 f \n")
sb.WriteString("0000000009 00000 n \n")
sb.WriteString("0000000058 00000 n \n")
sb.WriteString("0000000115 00000 n \n")
sb.WriteString("0000000241 00000 n \n")
sb.WriteString("0000000354 00000 n \n")
sb.WriteString("trailer << /Size 6 /Root 1 0 R >>\n")
sb.WriteString("startxref\n")
sb.WriteString("420\n")
sb.WriteString("%%EOF")
return sb.String()
}

func (p *PDFReport) SaveToFile(filepath string) bool {
content := p.GenerateMinimalPDF()
file, err := os.Create(filepath)
if err != nil {
return false
}
defer file.Close()
_, err = file.WriteString(content)
if err != nil {
return false
}
return true
}
