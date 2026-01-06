package plugins

import (
	"os"
	"path/filepath"
)

func Cursor(config PluginConfig) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	settingsPath := filepath.Join(
		home,
		"Library/Application Support/Cursor/User/settings.json",
	)

	// Use mode-specific settings from custom field if configured
	if settings := config.GetModeSettings(); len(settings) > 0 {
		return UpdateJSONSettings(settingsPath, settings)
	}

	// Fall back to legacy theme-only configuration
	theme := config.Night
	defaultTheme := "Default Dark+"

	if config.IsLight {
		theme = config.Day
		defaultTheme = "Default Light+"
	}

	if theme == "" {
		theme = defaultTheme
	}

	return UpdateJSONTheme(settingsPath, "workbench.colorTheme", theme)
}
