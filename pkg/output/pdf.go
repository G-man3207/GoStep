package output

import (
	"bytes"
	"fmt"
	"image/png"
	"time"

	"github.com/gustaf/go-test/pkg/recorder"
	"github.com/jung-kurt/gofpdf"
)

func SavePDF(steps []recorder.Step, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 15)

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 24)
	pdf.Cell(190, 10, "Step Recording")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 10, time.Now().Format("2006-01-02 15:04:05"))

	for i, step := range steps {
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(190, 10, fmt.Sprintf("Step %d", i+1))
		pdf.Ln(10)

		if step.Description != "" {
			pdf.SetFillColor(227, 242, 253) // Light blue background
			pdf.Rect(10, pdf.GetY(), 190, 10, "F")
			pdf.Cell(190, 10, fmt.Sprintf("Note: %s", step.Description))
			pdf.Ln(12)
		}

		var buf bytes.Buffer
		if err := png.Encode(&buf, step.Screenshot); err != nil {
			return fmt.Errorf("failed to encode screenshot: %w", err)
		}

		imgOpts := gofpdf.ImageOptions{
			ImageType: "PNG",
		}

		// Calculate image dimensions to fit page width while maintaining aspect ratio
		pageWidth := 190.0 // mm
		imgWidth := float64(step.Screenshot.Bounds().Dx())
		imgHeight := float64(step.Screenshot.Bounds().Dy())
		ratio := imgHeight / imgWidth
		width := pageWidth
		height := pageWidth * ratio

		pdf.RegisterImageOptionsReader(fmt.Sprintf("img%d", i), imgOpts, &buf)
		pdf.Image(fmt.Sprintf("img%d", i), 10, pdf.GetY(), width, height, false, "", 0, "")
	}

	return pdf.OutputFileAndClose(outputPath)
}
