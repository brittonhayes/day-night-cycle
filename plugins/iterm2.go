package plugins

import (
	"fmt"
	"os/exec"
)

func ITerm2(config PluginConfig) error {
	preset := config.Night
	if config.IsLight {
		preset = config.Day
	}

	if preset == "" {
		mode := "night"
		if config.IsLight {
			mode = "day"
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
