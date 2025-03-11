package output

import (
	"bytes"
	"fmt"
	"html/template"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/gustaf/go-test/pkg/recorder"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Step Recording - {{.Timestamp}}</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .step {
            background-color: white;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .step-header {
            margin-bottom: 10px;
        }
        .description {
            margin-top: 10px;
            padding: 10px;
            background-color: #e3f2fd;
            border-radius: 4px;
        }
        .screenshot {
            max-width: 100%;
            height: auto;
            border: 1px solid #ddd;
            border-radius: 4px;
            margin-top: 10px;
        }
    </style>
</head>
<body>
    <h1>Step Recording - {{.Timestamp}}</h1>
    {{range .Steps}}
    <div class="step">
        {{if .Description}}
        <div class="description">
            {{.Description}}
        </div>
        {{end}}
        <img class="screenshot" src="{{.ImagePath}}" alt="Screenshot">
    </div>
    {{end}}
</body>
</html>
`

type htmlStep struct {
	Timestamp   time.Time
	Action      string
	Description string
	Coordinates struct{ X, Y int }
	ImagePath   string
}

type htmlData struct {
	Timestamp time.Time
	Steps     []htmlStep
}

// SaveHTML saves the recording as an HTML file
func SaveHTML(steps []recorder.Step, outputPath string) error {
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	imagesDir := filepath.Join(outputDir, "images")
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return fmt.Errorf("failed to create images directory: %w", err)
	}

	data := htmlData{
		Timestamp: time.Now(),
		Steps:     make([]htmlStep, len(steps)),
	}

	for i, step := range steps {
		imgPath := filepath.Join("images", fmt.Sprintf("step_%d.png", i+1))
		imgFile, err := os.Create(filepath.Join(outputDir, imgPath))
		if err != nil {
			return fmt.Errorf("failed to create image file: %w", err)
		}
		if err := png.Encode(imgFile, step.Screenshot); err != nil {
			imgFile.Close()
			return fmt.Errorf("failed to encode image: %w", err)
		}
		imgFile.Close()

		data.Steps[i] = htmlStep{
			Timestamp:   step.Timestamp,
			Action:      step.Action,
			Description: step.Description,
			ImagePath:   imgPath,
		}
		data.Steps[i].Coordinates.X = step.Coordinates.X
		data.Steps[i].Coordinates.Y = step.Coordinates.Y
	}

	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}

	return nil
}
