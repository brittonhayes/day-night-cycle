package plugins

import (
	"os"
	"path/filepath"
)

func ClaudeCode(config PluginConfig) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	settingsPath := filepath.Join(home, ".claude/settings.json")

	// Use mode-specific settings from custom field if configured
	if settings := config.GetModeSettings(); len(settings) > 0 {
		return UpdateJSONSettings(settingsPath, settings)
	}

	// Fall back to legacy theme-only configuration
	theme := "dark"
	if config.IsLight {
		theme = "light"
	}

	return UpdateJSONTheme(settingsPath, "theme", theme)
}
