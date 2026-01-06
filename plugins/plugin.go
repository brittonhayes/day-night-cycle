package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PluginConfig provides theme configuration to plugins.
type PluginConfig struct {
	IsLight bool           // Whether to apply light mode
	Light   string         // Primary light mode value (theme/preset/colorscheme)
	Dark    string         // Primary dark mode value (theme/preset/colorscheme)
	Custom  map[string]any // Additional plugin-specific configuration
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
}

// NewPluginConfig creates a PluginConfig from plugin configuration.
func NewPluginConfig(pluginCfg map[string]interface{}, isLight bool) PluginConfig {
	config := PluginConfig{
		IsLight: isLight,
		Custom:  make(map[string]any),
	}

	lightKeys := []string{"light", "light_theme", "light_preset", "light_colorscheme"}
	darkKeys := []string{"dark", "dark_theme", "dark_preset", "dark_colorscheme"}

	extractedKeys := make(map[string]bool)

	for _, key := range lightKeys {
		if val, ok := pluginCfg[key].(string); ok {
			config.Light = val
			extractedKeys[key] = true
			break
		}
	}

	for _, key := range darkKeys {
		if val, ok := pluginCfg[key].(string); ok {
			config.Dark = val
			extractedKeys[key] = true
			break
		}
	}

	for k, v := range pluginCfg {
		if !extractedKeys[k] {
			config.Custom[k] = v
		}
	}

	return config
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
