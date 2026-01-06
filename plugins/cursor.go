package plugins

import (
	"os"
	"path/filepath"
)

// Cursor updates Cursor settings.json.
func Cursor(config PluginConfig) error {
	theme := config.Dark
	defaultTheme := "Default Dark+"

	if config.IsLight {
		theme = config.Light
		defaultTheme = "Default Light+"
	}

	// Use default if theme not configured
	if theme == "" {
		theme = defaultTheme
	}

	settingsPath := filepath.Join(
		os.Getenv("HOME"),
		"Library/Application Support/Cursor/User/settings.json",
	)

	return UpdateJSONTheme(settingsPath, "workbench.colorTheme", theme)
}
