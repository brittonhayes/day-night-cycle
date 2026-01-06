package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PluginConfig provides theme configuration to plugins.
// This is the source of truth for plugin configuration structure.
type PluginConfig struct {
	IsLight bool           `yaml:"-"`             // Whether to apply day mode (set at runtime)
	Day     string         `yaml:"day,omitempty"`  // Primary day mode value (theme/preset/colorscheme)
	Night   string         `yaml:"night,omitempty"` // Primary night mode value (theme/preset/colorscheme)
	Custom  map[string]any `yaml:"custom,omitempty"` // Additional plugin-specific configuration (supports "day" and "night" keys for mode-specific settings)
}

// Plugin is the signature for all plugin functions.
type Plugin func(config PluginConfig) error

// Registry holds all registered plugins.
var Registry = map[string]Plugin{
	"iterm2":       ITerm2,
	"cursor":       Cursor,
	"claude-code":  ClaudeCode,
	"neovim":       Neovim,
	"macos-system": MacOSSystem,
	"sublime":      Sublime,
	"pycharm":      PyCharm,
}

func UpdateJSONTheme(path, key, value string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
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

func UpdateJSONSettings(path string, updates map[string]any) error {
	if len(updates) == 0 {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	for key, value := range updates {
		settings[key] = value
	}

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

func ExpandPath(path string) (string, error) {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[1:]), nil
	}
	return path, nil
}

func (c PluginConfig) GetModeSettings() map[string]any {
	if len(c.Custom) == 0 {
		return nil
	}

	key := "night"
	if c.IsLight {
		key = "day"
	}

	settings, ok := c.Custom[key].(map[string]any)
	if !ok {
		return nil
	}

	return settings
}
