package plugins

import (
	"fmt"
	"os"
	"os/exec"
)

// MacOSSystem sets system-wide appearance.
func MacOSSystem(cfg map[string]interface{}, isLight bool) error {
	darkMode := "true"
	if isLight {
		darkMode = "false"
	}

	script := fmt.Sprintf(`
tell application "System Events"
	tell appearance preferences
		set dark mode to %s
	end tell
end tell
`, darkMode)

	cmd := exec.Command("osascript", "-e", script)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("osascript failed: %w: %s", err, output)
	}

	// Optional wallpaper support
	wallpaperKey := "dark_wallpaper"
	if isLight {
		wallpaperKey = "light_wallpaper"
	}

	if wallpaper, ok := cfg[wallpaperKey].(string); ok {
		fullPath := ExpandPath(wallpaper)

		if _, err := os.Stat(fullPath); err != nil {
			fmt.Printf("    Warning: wallpaper file not found: %s\n", fullPath)
			return nil
		}

		wallpaperScript := fmt.Sprintf(`
tell application "Finder"
	set desktop picture to POSIX file "%s"
end tell
`, fullPath)

		cmd := exec.Command("osascript", "-e", wallpaperScript)
		_ = cmd.Run()
	}

	return nil
}
