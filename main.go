package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the YAML configuration
type Config struct {
	Location LocationConfig `yaml:"location"`
	Plugins  []PluginConfig `yaml:"plugins"`
}

// LocationConfig holds geographic location settings
type LocationConfig struct {
	Name      string  `yaml:"name"`
	Latitude  float64 `yaml:"latitude"`
	Longitude float64 `yaml:"longitude"`
	Timezone  string  `yaml:"timezone"`
}

// PluginConfig holds configuration for a single plugin
type PluginConfig struct {
	Name    string                 `yaml:"name"`
	Enabled bool                   `yaml:"enabled"`
	Config  map[string]interface{} `yaml:",inline"`
}

// PluginFunc is the signature for all plugin functions
type PluginFunc func(pluginCfg map[string]interface{}, isLight bool) error

var Version = "dev"

func main() {
	configPath := flag.String("config", defaultConfigPath(), "path to config file")
	flag.Usage = printUsage
	flag.Parse()

	if flag.NArg() < 1 {
		printUsage()
		os.Exit(1)
	}

	command := flag.Arg(0)

	switch command {
	case "auto":
		runAuto(*configPath)
	case "light":
		runMode(*configPath, true)
	case "dark":
		runMode(*configPath, false)
	case "status":
		runStatus(*configPath)
	case "next":
		runNext(*configPath)
	case "schedule":
		runSchedule(*configPath)
	case "version":
		fmt.Printf("day-night-cycle version %s\n", Version)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `day-night-cycle - Automatically switch themes based on sunrise/sunset

Usage:
  day-night-cycle [flags] <command>

Commands:
  auto      Apply mode based on current time
  light     Force light mode
  dark      Force dark mode
  status    Show current status and schedule
  next      Show next transition time
  schedule  Generate launchd schedule
  version   Show version

Flags:
`)
	flag.PrintDefaults()
}

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "config.yaml"
	}
	return filepath.Join(home, ".config", "day-night-cycle", "config.yaml")
}

func loadConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config: %v\n", err)
		os.Exit(1)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing config: %v\n", err)
		os.Exit(1)
	}

	return cfg
}

func loadLocation(tz string) *time.Location {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading timezone %s: %v\n", tz, err)
		os.Exit(1)
	}
	return loc
}

func runAuto(configPath string) {
	cfg := loadConfig(configPath)
	loc := loadLocation(cfg.Location.Timezone)

	sunrise, sunset := calculateSolarTimes(
		cfg.Location.Latitude,
		cfg.Location.Longitude,
		time.Now().In(loc),
	)

	now := time.Now().In(loc)
	isLight := now.After(sunrise) && now.Before(sunset)

	applyMode(cfg, isLight)
}

func runMode(configPath string, isLight bool) {
	cfg := loadConfig(configPath)
	applyMode(cfg, isLight)
}

func applyMode(cfg Config, isLight bool) {
	mode := "dark"
	if isLight {
		mode = "light"
	}
	fmt.Printf("\nApplying %s mode...\n", mode)

	success := 0
	total := 0

	for _, pluginCfg := range cfg.Plugins {
		if !pluginCfg.Enabled {
			continue
		}

		pluginFunc, exists := plugins[pluginCfg.Name]
		if !exists {
			fmt.Printf("  ✗ %s: unknown plugin\n", pluginCfg.Name)
			continue
		}

		total++
		err := pluginFunc(pluginCfg.Config, isLight)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", pluginCfg.Name, err)
		} else {
			fmt.Printf("  ✓ %s\n", pluginCfg.Name)
			success++
		}
	}

	fmt.Printf("\nCompleted: %d/%d plugins successful\n", success, total)
}

func runStatus(configPath string) {
	cfg := loadConfig(configPath)
	loc := loadLocation(cfg.Location.Timezone)

	now := time.Now().In(loc)
	sunrise, sunset := calculateSolarTimes(
		cfg.Location.Latitude,
		cfg.Location.Longitude,
		now,
	)

	isLight := now.After(sunrise) && now.Before(sunset)
	currentMode := "dark"
	if isLight {
		currentMode = "light"
	}

	fmt.Printf("\nCurrent mode: %s\n", currentMode)
	fmt.Printf("Today's sunrise: %s\n", sunrise.Format("3:04 PM"))
	fmt.Printf("Today's sunset: %s\n", sunset.Format("3:04 PM"))

	var next time.Time
	var kind string
	if now.Before(sunrise) {
		next = sunrise
		kind = "sunrise"
	} else if now.Before(sunset) {
		next = sunset
		kind = "sunset"
	} else {
		tomorrow := now.Add(24 * time.Hour)
		next, _ = calculateSolarTimes(cfg.Location.Latitude, cfg.Location.Longitude, tomorrow)
		kind = "sunrise"
	}
	fmt.Printf("Next transition: %s (%s)\n", next.Format("3:04 PM"), kind)

	fmt.Println("\nConfigured plugins:")
	for _, pluginCfg := range cfg.Plugins {
		if pluginCfg.Enabled {
			fmt.Printf("  • %s\n", pluginCfg.Name)
		}
	}
	fmt.Println()
}

func runNext(configPath string) {
	cfg := loadConfig(configPath)
	loc := loadLocation(cfg.Location.Timezone)

	now := time.Now().In(loc)
	sunrise, sunset := calculateSolarTimes(
		cfg.Location.Latitude,
		cfg.Location.Longitude,
		now,
	)

	var next time.Time
	var kind string
	if now.Before(sunrise) {
		next = sunrise
		kind = "sunrise"
	} else if now.Before(sunset) {
		next = sunset
		kind = "sunset"
	} else {
		tomorrow := now.Add(24 * time.Hour)
		next, _ = calculateSolarTimes(cfg.Location.Latitude, cfg.Location.Longitude, tomorrow)
		kind = "sunrise"
	}

	fmt.Printf("Next transition: %s (%s)\n", next.Format("3:04 PM"), kind)
}
