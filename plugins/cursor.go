package plugins

import (
	"os"
	"path/filepath"
)

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

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	settingsPath := filepath.Join(
		home,
		"Library/Application Support/Cursor/User/settings.json",
	)

	return UpdateJSONTheme(settingsPath, "workbench.colorTheme", theme)
}
