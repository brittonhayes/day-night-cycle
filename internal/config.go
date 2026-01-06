package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the YAML configuration.
type Config struct {
	Location LocationConfig `yaml:"location"`
	Plugins  []PluginConfig `yaml:"plugins"`
}

// LocationConfig holds geographic location settings.
type LocationConfig struct {
	Latitude  float64 `yaml:"latitude"`
	Longitude float64 `yaml:"longitude"`
	Timezone  string  `yaml:"timezone"`
}

// PluginConfig holds configuration for a single plugin.
type PluginConfig struct {
	Name    string                 `yaml:"name"`
	Enabled bool                   `yaml:"enabled"`
	Config  map[string]interface{} `yaml:",inline"`
}

// DefaultPath returns the default configuration file path.
func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "config.yaml"
	}
	return filepath.Join(home, ".config", "day-night-cycle", "config.yaml")
}

// Load reads and parses the configuration file.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

// LoadLocation loads the timezone location.
func LoadLocation(tz string) (*time.Location, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("loading timezone %s: %w", tz, err)
	}
	return loc, nil
}
