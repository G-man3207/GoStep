//go:build windows
// +build windows

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/gustaf/go-test/pkg/config"
)

func ShowSettingsDialog(window fyne.Window, settings *config.Settings) error {
	form := container.NewVBox(
		widget.NewLabel("Output Settings"),
		widget.NewLabel("Output directory: "+settings.OutputDir),
	)

	dlg := dialog.NewCustom("Settings", "Close", form, window)
	dlg.Resize(fyne.NewSize(300, 150))
	dlg.Show()

	return nil
}
