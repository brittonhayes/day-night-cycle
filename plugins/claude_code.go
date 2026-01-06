package plugins

import (
	"os"
	"path/filepath"
)

// ClaudeCode updates Claude Code settings.json.
func ClaudeCode(cfg map[string]interface{}, isLight bool) error {
	theme := "dark"
	if isLight {
		theme = "light"
	}

	settingsPath := filepath.Join(os.Getenv("HOME"), ".claude/settings.json")
	return UpdateJSONTheme(settingsPath, "theme", theme)
}
