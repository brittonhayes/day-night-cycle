package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.daynightcycle.schedule</string>
	<key>ProgramArguments</key>
	<array>
		<string>{{.BinaryPath}}</string>
		<string>--config</string>
		<string>{{.ConfigPath}}</string>
		<string>auto</string>
	</array>
	<key>StartCalendarInterval</key>
	<array>
		<dict>
			<key>Hour</key>
			<integer>{{.SunriseHour}}</integer>
			<key>Minute</key>
			<integer>{{.SunriseMinute}}</integer>
		</dict>
		<dict>
			<key>Hour</key>
			<integer>{{.SunsetHour}}</integer>
			<key>Minute</key>
			<integer>{{.SunsetMinute}}</integer>
		</dict>
	</array>
	<key>StandardOutPath</key>
	<string>{{.LogPath}}/schedule.log</string>
	<key>StandardErrorPath</key>
	<string>{{.LogPath}}/schedule.error.log</string>
</dict>
</plist>`

func runSchedule(configPath string) {
	cfg := loadConfig(configPath)
	loc := loadLocation(cfg.Location.Timezone)

	now := time.Now().In(loc)
	sunrise, sunset := calculateSolarTimes(
		cfg.Location.Latitude,
		cfg.Location.Longitude,
		now,
	)

	// Get binary path
	binaryPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting executable path: %v\n", err)
		os.Exit(1)
	}

	// Resolve symlinks
	binaryPath, err = filepath.EvalSymlinks(binaryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving symlinks: %v\n", err)
		os.Exit(1)
	}

	// Get absolute config path
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		absConfigPath = configPath
	}

	// Set up paths
	home, _ := os.UserHomeDir()
	launchdDir := filepath.Join(home, "Library/LaunchAgents")
	plistPath := filepath.Join(launchdDir, "com.daynightcycle.schedule.plist")
	logPath := filepath.Join(filepath.Dir(absConfigPath), "logs")

	// Ensure directories exist
	if err := os.MkdirAll(launchdDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating LaunchAgents directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(logPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating logs directory: %v\n", err)
		os.Exit(1)
	}

	// Prepare template data
	data := map[string]interface{}{
		"BinaryPath":    binaryPath,
		"ConfigPath":    absConfigPath,
		"SunriseHour":   sunrise.Hour(),
		"SunriseMinute": sunrise.Minute(),
		"SunsetHour":    sunset.Hour(),
		"SunsetMinute":  sunset.Minute(),
		"LogPath":       logPath,
	}

	// Parse and execute template
	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing template: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create(plistPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating plist file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		fmt.Fprintf(os.Stderr, "error writing plist: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nLaunchd schedule created successfully\n")
	fmt.Printf("\nSchedule for %s:\n", now.Format("Monday, January 2, 2006"))
	fmt.Printf("  Sunrise: %s\n", sunrise.Format("3:04 PM"))
	fmt.Printf("  Sunset:  %s\n", sunset.Format("3:04 PM"))
	fmt.Printf("\nPlist file: %s\n", plistPath)
	fmt.Printf("Logs directory: %s\n", logPath)
	fmt.Printf("\nTo enable automatic theme switching:\n")
	fmt.Printf("  launchctl unload %s 2>/dev/null || true\n", plistPath)
	fmt.Printf("  launchctl load %s\n", plistPath)
	fmt.Printf("\nTo disable automatic theme switching:\n")
	fmt.Printf("  launchctl unload %s\n", plistPath)
	fmt.Println()
}
