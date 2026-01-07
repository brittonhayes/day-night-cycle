package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/brittonhayes/day-night-cycle/internal"
	"github.com/brittonhayes/day-night-cycle/plugins"
)

var Version = "dev"

func main() {
	configPath := flag.String("config", internal.DefaultPath(), "path to config file")
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

func runAuto(configPath string) {
	cfg, err := internal.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	loc, err := internal.LoadLocation(cfg.Location.Timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	sunrise, sunset := internal.CalculateTimes(
		cfg.Location.Latitude,
		cfg.Location.Longitude,
		time.Now().In(loc),
	)

	sunrise, sunset = cfg.Location.ApplyOffsets(sunrise, sunset)

	now := time.Now().In(loc)
	isLight := now.After(sunrise) && now.Before(sunset)

	applyMode(cfg, isLight)
}

func runMode(configPath string, isLight bool) {
	cfg, err := internal.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	applyMode(cfg, isLight)
}

func applyMode(cfg internal.Config, isLight bool) {
	mode := "dark"
	if isLight {
		mode = "light"
	}
	fmt.Printf("\nApplying %s mode...\n", mode)

	success := 0
	total := 0

	for _, pluginEntry := range cfg.Plugins {
		if !pluginEntry.Enabled {
			continue
		}

		pluginFunc, exists := plugins.Registry[pluginEntry.Name]
		if !exists {
			fmt.Printf("  ✗ %s: unknown plugin\n", pluginEntry.Name)
			continue
		}

		total++
		config := pluginEntry.PluginConfig
		config.IsLight = isLight
		err := pluginFunc(config)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", pluginEntry.Name, err)
		} else {
			fmt.Printf("  ✓ %s\n", pluginEntry.Name)
			success++
		}
	}

	fmt.Printf("\nCompleted: %d/%d plugins successful\n", success, total)
}

func nextTransition(now, sunrise, sunset time.Time, loc internal.LocationConfig) (next time.Time, kind string) {
	if now.Before(sunrise) {
		return sunrise, "sunrise"
	}
	if now.Before(sunset) {
		return sunset, "sunset"
	}
	tomorrow := now.Add(24 * time.Hour)
	next, _ = internal.CalculateTimes(loc.Latitude, loc.Longitude, tomorrow)
	next, _ = loc.ApplyOffsets(next, time.Time{})
	return next, "sunrise"
}

func runStatus(configPath string) {
	cfg, err := internal.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	loc, err := internal.LoadLocation(cfg.Location.Timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	now := time.Now().In(loc)
	sunrise, sunset := internal.CalculateTimes(
		cfg.Location.Latitude,
		cfg.Location.Longitude,
		now,
	)

	sunrise, sunset = cfg.Location.ApplyOffsets(sunrise, sunset)

	isLight := now.After(sunrise) && now.Before(sunset)
	currentMode := "dark"
	if isLight {
		currentMode = "light"
	}

	fmt.Printf("\nCurrent mode: %s\n", currentMode)

	if cfg.Location.DayOffset != "" {
		fmt.Printf("Sunrise: %s (offset: %s)\n", sunrise.Format("3:04 PM"), cfg.Location.DayOffset)
	} else {
		fmt.Printf("Sunrise: %s\n", sunrise.Format("3:04 PM"))
	}

	if cfg.Location.NightOffset != "" {
		fmt.Printf("Sunset: %s (offset: %s)\n", sunset.Format("3:04 PM"), cfg.Location.NightOffset)
	} else {
		fmt.Printf("Sunset: %s\n", sunset.Format("3:04 PM"))
	}

	next, kind := nextTransition(now, sunrise, sunset, cfg.Location)
	fmt.Printf("Next transition: %s (%s)\n", next.Format("3:04 PM"), kind)

	fmt.Println("\nConfigured plugins:")
	for _, pluginEntry := range cfg.Plugins {
		if pluginEntry.Enabled {
			fmt.Printf("  • %s\n", pluginEntry.Name)
		}
	}
	fmt.Println()
}

func runNext(configPath string) {
	cfg, err := internal.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	loc, err := internal.LoadLocation(cfg.Location.Timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	now := time.Now().In(loc)
	sunrise, sunset := internal.CalculateTimes(
		cfg.Location.Latitude,
		cfg.Location.Longitude,
		now,
	)

	sunrise, sunset = cfg.Location.ApplyOffsets(sunrise, sunset)

	next, kind := nextTransition(now, sunrise, sunset, cfg.Location)
	fmt.Printf("Next transition: %s (%s)\n", next.Format("3:04 PM"), kind)
}

func runSchedule(configPath string) {
	cfg, err := internal.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	loc, err := internal.LoadLocation(cfg.Location.Timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	now := time.Now().In(loc)
	sunrise, sunset := internal.CalculateTimes(
		cfg.Location.Latitude,
		cfg.Location.Longitude,
		now,
	)

	sunrise, sunset = cfg.Location.ApplyOffsets(sunrise, sunset)

	if err := internal.Generate(configPath, sunrise, sunset); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
