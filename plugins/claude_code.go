package plugins

import (
	"os"
	"path/filepath"
)

// ClaudeCode updates Claude Code settings.json.
func ClaudeCode(config PluginConfig) error {
	theme := "dark"
	if config.IsLight {
		theme = "light"
	}

	settingsPath := filepath.Join(os.Getenv("HOME"), ".claude/settings.json")
	return UpdateJSONTheme(settingsPath, "theme", theme)
}
