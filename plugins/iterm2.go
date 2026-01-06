package plugins

import (
	"fmt"
	"os/exec"
)

// ITerm2 switches iTerm2 color presets via AppleScript.
func ITerm2(config PluginConfig) error {
	preset := config.Dark
	if config.IsLight {
		preset = config.Light
	}

	if preset == "" {
		mode := "dark"
		if config.IsLight {
			mode = "light"
		}
		return fmt.Errorf("missing %s preset configuration", mode)
	}

	script := fmt.Sprintf(`
tell application "iTerm"
	repeat with aWindow in windows
		repeat with aTab in tabs of aWindow
			repeat with aSession in sessions of aTab
				tell aSession
					set color preset to "%s"
				end tell
			end repeat
		end repeat
	end repeat
end tell
`, preset)

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("osascript failed: %w: %s", err, output)
	}

	return nil
}
