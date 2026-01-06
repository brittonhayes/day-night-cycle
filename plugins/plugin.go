package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Func is the signature for all plugin functions.
type Func func(pluginCfg map[string]interface{}, isLight bool) error

// Registry holds all registered plugins.
var Registry = map[string]Func{
	"iterm2":       ITerm2,
	"cursor":       Cursor,
	"claude-code":  ClaudeCode,
	"neovim":       Neovim,
	"macos-system": MacOSSystem,
}

// UpdateJSONTheme updates a JSON file with a new theme value.
func UpdateJSONTheme(path, key, value string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	// Check if already set (optimization)
	if current, ok := settings[key].(string); ok && current == value {
		return nil
	}

	settings[key] = value

	output, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, output, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}

	return nil
}

// ExpandPath handles ~ expansion in file paths.
func ExpandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}
