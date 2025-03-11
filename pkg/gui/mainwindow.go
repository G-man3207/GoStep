//go:build windows
// +build windows

package gui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/gustaf/go-test/pkg/config"
	"github.com/gustaf/go-test/pkg/output"
	"github.com/gustaf/go-test/pkg/recorder"
)

// RecorderWindow represents the main application window
type RecorderWindow struct {
	window     fyne.Window
	app        fyne.App
	recorder   *recorder.Recorder
	settings   *config.Settings
	isRunning  bool
	outputPath string
}

// NewRecorderWindow creates and shows the main application window
func NewRecorderWindow(settings *config.Settings) error {
	app := app.New()
	window := app.NewWindow("GoStep")

	rw := &RecorderWindow{
		window:    window,
		app:       app,
		settings:  settings,
		recorder:  recorder.NewRecorder(),
		isRunning: false,
	}

	recordBtn := widget.NewButton("Start Recording", nil)
	status := widget.NewLabel("Ready")
	outputFormatSelect := widget.NewSelect([]string{"HTML", "PDF"}, nil)
	outputFormatSelect.SetSelected(settings.OutputFormat)

	outputLocationLabel := widget.NewLabel(fmt.Sprintf("Output Location: %s", settings.OutputDir))
	outputLocationLabel.Wrapping = fyne.TextWrapBreak

	settingsBtn := widget.NewButton("Settings", func() {
		ShowSettingsDialog(window, settings)
	})

	toolbar := container.NewHBox(
		recordBtn,
		layout.NewSpacer(),
		widget.NewLabel("Format:"),
		outputFormatSelect,
		settingsBtn,
	)

	statusArea := container.NewVBox(
		container.NewHBox(
			widget.NewIcon(theme.InfoIcon()),
			status,
		),
		outputLocationLabel,
	)

	content := container.NewVBox(
		toolbar,
		widget.NewSeparator(),
		statusArea,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(600, 400))

	outputFormatSelect.OnChanged = func(format string) {
		settings.OutputFormat = format
		rw.updateOutputPath()
	}

	recordBtn.OnTapped = func() {
		if !rw.isRunning {
			if err := rw.startRecording(); err != nil {
				dialog.ShowError(err, window)
				return
			}
			recordBtn.SetText("Stop Recording")
			status.SetText("Recording... Click anywhere to capture.")
			rw.isRunning = true
		} else {
			if err := rw.stopRecording(); err != nil {
				dialog.ShowError(err, window)
				return
			}
			recordBtn.SetText("Start Recording")
			status.SetText("Ready")
			rw.isRunning = false
		}
	}

	rw.updateOutputPath()

	window.Show()
	app.Run()

	return nil
}

func (rw *RecorderWindow) updateOutputPath() {
	timestamp := time.Now().Format("2006-01-02_150405")
	ext := "html"
	if rw.settings.OutputFormat == "PDF" {
		ext = "pdf"
	}
	filename := fmt.Sprintf("recording_%s.%s", timestamp, ext)
	rw.outputPath = filepath.Join(rw.settings.OutputDir, filename)
}

func (rw *RecorderWindow) startRecording() error {
	log.Printf("Starting recording...")
	if err := rw.recorder.Start(); err != nil {
		log.Printf("Failed to start recording: %v", err)
		return err
	}
	log.Printf("Recording started successfully")
	return nil
}

func (rw *RecorderWindow) stopRecording() error {
	log.Printf("Stopping recording...")
	if err := rw.recorder.Stop(); err != nil {
		log.Printf("Failed to stop recording: %v", err)
		return err
	}

	steps := rw.recorder.GetSteps()
	log.Printf("Retrieved %d steps from recorder", len(steps))
	if len(steps) == 0 {
		log.Printf("No steps were recorded")
		return fmt.Errorf("no steps recorded")
	}

	rw.showImageEditor(steps)
	return nil
}

func (rw *RecorderWindow) showImageEditor(steps []recorder.Step) {
	previewWindow := rw.app.NewWindow("Preview Recording")
	previewWindow.Resize(fyne.NewSize(800, 600))

	vbox := container.NewVBox()

	scrollContainer := container.NewVScroll(vbox)
	scrollContainer.SetMinSize(fyne.NewSize(780, 500))

	wrapper := container.NewPadded(scrollContainer)

	stepWidgets := make([]*fyne.Container, len(steps))

	saveBtn := widget.NewButton("Save", func() {
		if err := os.MkdirAll(filepath.Dir(rw.outputPath), 0755); err != nil {
			dialog.ShowError(fmt.Errorf("failed to create output directory: %w", err), previewWindow)
			return
		}

		switch rw.settings.OutputFormat {
		case "HTML":
			if err := output.SaveHTML(steps, rw.outputPath); err != nil {
				dialog.ShowError(fmt.Errorf("failed to save HTML: %w", err), previewWindow)
				return
			}
		case "PDF":
			if err := output.SavePDF(steps, rw.outputPath); err != nil {
				dialog.ShowError(fmt.Errorf("failed to save PDF: %w", err), previewWindow)
				return
			}
		default:
			dialog.ShowError(fmt.Errorf("unsupported output format: %s", rw.settings.OutputFormat), previewWindow)
			return
		}
		previewWindow.Close()
	})

	toolbar := container.NewHBox(
		saveBtn,
		layout.NewSpacer(),
		widget.NewLabelWithStyle(
			fmt.Sprintf("Total Steps: %d", len(steps)),
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),
	)

	for i, step := range steps {
		img := canvas.NewImageFromImage(step.Screenshot)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(400, 300))

		imgContainer := container.NewHBox(layout.NewSpacer(), img, layout.NewSpacer())
		imgContainer.Resize(fyne.NewSize(750, 525))

		buttonSize := fyne.NewSize(100, 32)

		deleteBtn := widget.NewButton("Delete", nil)
		deleteBtn.Resize(buttonSize)

		moveUpBtn := widget.NewButton("Move Up", nil)
		moveUpBtn.Resize(buttonSize)

		moveDownBtn := widget.NewButton("Move Down", nil)
		moveDownBtn.Resize(buttonSize)

		descLabel := widget.NewTextGrid()
		descLabel.SetText(step.Description)

		editBtn := widget.NewButton("Edit Description", nil)
		editBtn.OnTapped = func() {
			entry := widget.NewMultiLineEntry()
			entry.SetText(descLabel.Text())
			entry.SetPlaceHolder("Add description...")
			entry.Wrapping = fyne.TextWrapWord
			entry.Resize(fyne.NewSize(600, 400))

			entryContainer := container.NewPadded(
				container.NewScroll(entry),
			)
			entryContainer.Resize(fyne.NewSize(600, 400))

			dialog := dialog.NewCustomConfirm("Edit Description", "Save", "Cancel",
				entryContainer,
				func(save bool) {
					if save {
						descLabel.SetText(entry.Text)
						steps[i].Description = entry.Text
					}
				},
				previewWindow,
			)
			dialog.Resize(fyne.NewSize(700, 500))
			dialog.Show()
		}

		descContainer := container.NewVBox(
			descLabel,
		)

		controls := container.NewHBox(
			deleteBtn,
			moveUpBtn,
			moveDownBtn,
			editBtn,
		)

		stepContainer := container.NewVBox(
			widget.NewLabelWithStyle(
				fmt.Sprintf("Step %d", i+1),
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			),
			container.NewPadded(
				container.NewVBox(
					imgContainer,
					container.NewPadded(descContainer),
					controls,
				),
			),
			widget.NewSeparator(),
		)
		stepWidgets[i] = stepContainer

		deleteBtn.OnTapped = func() {
			vbox.Remove(stepContainer)
			newSteps := make([]recorder.Step, 0)
			for _, s := range steps {
				if s != step {
					newSteps = append(newSteps, s)
				}
			}
			steps = newSteps
		}

		moveUpBtn.OnTapped = func() {
			if i > 0 {
				steps[i], steps[i-1] = steps[i-1], steps[i]
				stepWidgets[i], stepWidgets[i-1] = stepWidgets[i-1], stepWidgets[i]

				vbox.RemoveAll()
				for _, widget := range stepWidgets {
					vbox.Add(widget)
				}
			}
		}

		moveDownBtn.OnTapped = func() {
			if i < len(steps)-1 {
				steps[i], steps[i+1] = steps[i+1], steps[i]
				stepWidgets[i], stepWidgets[i+1] = stepWidgets[i+1], stepWidgets[i]

				vbox.RemoveAll()
				for _, widget := range stepWidgets {
					vbox.Add(widget)
				}
			}
		}

		vbox.Add(stepContainer)
	}

	mainContainer := container.NewBorder(
		toolbar, nil, nil, nil,
		wrapper,
	)

	previewWindow.SetContent(mainContainer)
	previewWindow.Show()
}
