package plugins

import (
	"os"
	"path/filepath"
)

// Cursor updates Cursor settings.json.
func Cursor(cfg map[string]interface{}, isLight bool) error {
	themeKey := "dark_theme"
	defaultTheme := "Default Dark+"
	if isLight {
		themeKey = "light_theme"
		defaultTheme = "Default Light+"
	}

	theme, ok := cfg[themeKey].(string)
	if !ok {
		theme = defaultTheme
	}

	settingsPath := filepath.Join(
		os.Getenv("HOME"),
		"Library/Application Support/Cursor/User/settings.json",
	)

	return UpdateJSONTheme(settingsPath, "workbench.colorTheme", theme)
}
