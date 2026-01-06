package internal

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

// Generate creates a launchd plist file for automatic scheduling.
func Generate(configPath string, sunrise, sunset time.Time) error {
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}

	binaryPath, err = filepath.EvalSymlinks(binaryPath)
	if err != nil {
		return fmt.Errorf("resolving symlinks: %w", err)
	}

	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		absConfigPath = configPath
	}

	home, _ := os.UserHomeDir()
	launchdDir := filepath.Join(home, "Library/LaunchAgents")
	plistPath := filepath.Join(launchdDir, "com.daynightcycle.schedule.plist")
	logPath := filepath.Join(filepath.Dir(absConfigPath), "logs")

	if err := os.MkdirAll(launchdDir, 0755); err != nil {
		return fmt.Errorf("creating LaunchAgents directory: %w", err)
	}

	if err := os.MkdirAll(logPath, 0755); err != nil {
		return fmt.Errorf("creating logs directory: %w", err)
	}

	data := map[string]interface{}{
		"BinaryPath":    binaryPath,
		"ConfigPath":    absConfigPath,
		"SunriseHour":   sunrise.Hour(),
		"SunriseMinute": sunrise.Minute(),
		"SunsetHour":    sunset.Hour(),
		"SunsetMinute":  sunset.Minute(),
		"LogPath":       logPath,
	}

	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	f, err := os.Create(plistPath)
	if err != nil {
		return fmt.Errorf("creating plist file: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("writing plist: %w", err)
	}

	fmt.Printf("\nLaunchd schedule created successfully\n")
	fmt.Printf("\nSchedule for %s:\n", time.Now().Format("Monday, January 2, 2006"))
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

	return nil
}
