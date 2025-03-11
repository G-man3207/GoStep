package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	OutputFormat     string `json:"output_format"` // "html" or "pdf"
	OutputDir        string `json:"output_dir"`
	MainWindowWidth  int    `json:"main_window_width"`
	MainWindowHeight int    `json:"main_window_height"`
}

func DefaultSettings() *Settings {
	return &Settings{
		OutputFormat:     "html",
		OutputDir:        filepath.Join(os.Getenv("USERPROFILE"), "Documents", "GoStep"),
		MainWindowWidth:  800,
		MainWindowHeight: 600,
	}
}

func LoadSettings(configPath string) (*Settings, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultSettings(), nil
		}
		return nil, err
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func (s *Settings) SaveSettings(configPath string) error {
	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
