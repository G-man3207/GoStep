package output

import (
	"bytes"
	"fmt"
	"html/template"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type Step struct {
	Screenshot  []byte
	Description string
	Timestamp   time.Time
	Action      string
}

func SaveAsHTML(steps []Step, outputPath string) error {
	const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Step Recording Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .step { margin-bottom: 30px; border-bottom: 1px solid #ccc; padding-bottom: 20px; }
        .timestamp { color: #666; font-size: 0.9em; }
        .description { margin: 10px 0; }
        img { max-width: 100%; border: 1px solid #ddd; margin-top: 10px; }
    </style>
</head>
<body>
    <h1>Step Recording Report</h1>
    {{range .Steps}}
    <div class="step">
        <div class="timestamp">{{.Timestamp.Format "2006-01-02 15:04:05"}}</div>
        <div class="description">{{.Description}}</div>
        <img src="data:image/png;base64,{{.Screenshot}}" alt="Screenshot">
    </div>
    {{end}}
</body>
</html>`

	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	data := struct {
		Steps []Step
	}{
		Steps: steps,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}

func SaveAsPDF(steps []Step, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")

	for i, step := range steps {
		pdf.AddPage()

		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 10, fmt.Sprintf("Step %d - %s", i+1, step.Timestamp.Format("2006-01-02 15:04:05")))
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 10, step.Description, "", "", false)
		pdf.Ln(5)

		imgReader := bytes.NewReader(step.Screenshot)
		img, err := png.Decode(imgReader)
		if err != nil {
			return fmt.Errorf("failed to decode image for step %d: %v", i+1, err)
		}

		tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("step_%d.png", i))
		f, err := os.Create(tmpFile)
		if err != nil {
			return fmt.Errorf("failed to create temporary file: %v", err)
		}

		if err := png.Encode(f, img); err != nil {
			f.Close()
			os.Remove(tmpFile)
			return fmt.Errorf("failed to encode image: %v", err)
		}
		f.Close()

		pdf.Image(tmpFile, 10, pdf.GetY(), 190, 0, false, "", 0, "")
		os.Remove(tmpFile)
	}

	return pdf.OutputFileAndClose(outputPath)
}
