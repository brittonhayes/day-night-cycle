package plugins

import (
	"os"
	"path/filepath"
)

func Sublime(config PluginConfig) error {
	colorScheme := config.Night
	defaultScheme := "Monokai.sublime-color-scheme"

	if config.IsLight {
		colorScheme = config.Day
		defaultScheme = "Breakers.sublime-color-scheme"
	}

	if colorScheme == "" {
		colorScheme = defaultScheme
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Try Sublime Text 4 first, then fall back to Sublime Text 3
	paths := []string{
		filepath.Join(home, "Library/Application Support/Sublime Text/Packages/User/Preferences.sublime-settings"),
		filepath.Join(home, "Library/Application Support/Sublime Text 4/Packages/User/Preferences.sublime-settings"),
		filepath.Join(home, "Library/Application Support/Sublime Text 3/Packages/User/Preferences.sublime-settings"),
	}

	var lastErr error
	for _, settingsPath := range paths {
		if _, err := os.Stat(settingsPath); err == nil {
			return UpdateJSONTheme(settingsPath, "color_scheme", colorScheme)
		}
		lastErr = err
	}

	return lastErr
}
