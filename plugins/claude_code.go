package plugins

import (
	"os"
	"path/filepath"
)

func ClaudeCode(config PluginConfig) error {
	theme := "dark"
	if config.IsLight {
		theme = "light"
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	settingsPath := filepath.Join(home, ".claude/settings.json")
	return UpdateJSONTheme(settingsPath, "theme", theme)
}
