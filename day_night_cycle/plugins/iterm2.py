"""iTerm2 plugin for day/night cycle automation."""

import subprocess
from typing import Optional
from .base import Plugin


class iTerm2Plugin(Plugin):
    """Plugin to control iTerm2 color presets."""

    @property
    def name(self) -> str:
        return "iterm2"

    def validate_config(self) -> tuple[bool, Optional[str]]:
        """Validate that required presets are configured."""
        light_preset = self.config.get('light_preset')
        dark_preset = self.config.get('dark_preset')

        if not light_preset:
            return False, "light_preset is required"
        if not dark_preset:
            return False, "dark_preset is required"

        return True, None

    def _set_preset(self, preset_name: str) -> bool:
        """
        Set iTerm2 color preset using AppleScript.

        Args:
            preset_name: Name of the color preset to activate

        Returns:
            True if successful, False otherwise
        """
        # Apply to all sessions in all windows to ensure it takes effect
        applescript = f'''
        tell application "iTerm"
            repeat with aWindow in windows
                repeat with aTab in tabs of aWindow
                    repeat with aSession in sessions of aTab
                        tell aSession
                            set color preset to "{preset_name}"
                        end tell
                    end repeat
                end repeat
            end repeat
        end tell
        '''

        try:
            result = subprocess.run(
                ['osascript', '-e', applescript],
                check=True,
                capture_output=True,
                text=True,
                timeout=5
            )
            # Check if there was any error output
            if result.stderr:
                print(f"    Warning: {result.stderr.strip()}")
                return False
            return True
        except subprocess.CalledProcessError as e:
            error_msg = e.stderr.strip() if e.stderr else str(e)
            print(f"    Error: {error_msg}")
            return False
        except subprocess.TimeoutExpired:
            print(f"    Error: Command timed out")
            return False
        except Exception as e:
            print(f"    Error: {e}")
            return False

    def set_light_mode(self) -> bool:
        """Set iTerm2 to light mode preset."""
        preset = self.config.get('light_preset', 'github-light')
        return self._set_preset(preset)

    def set_dark_mode(self) -> bool:
        """Set iTerm2 to dark mode preset."""
        preset = self.config.get('dark_preset', 'githubdark')
        return self._set_preset(preset)
