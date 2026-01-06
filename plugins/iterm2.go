package plugins

import (
	"fmt"
	"os/exec"
)

// ITerm2 switches iTerm2 color presets via AppleScript.
func ITerm2(cfg map[string]interface{}, isLight bool) error {
	presetKey := "dark_preset"
	if isLight {
		presetKey = "light_preset"
	}

	preset, ok := cfg[presetKey].(string)
	if !ok {
		return fmt.Errorf("missing %s configuration", presetKey)
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
