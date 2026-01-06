"""macOS System plugin for day/night cycle automation."""

import subprocess
import platform
from pathlib import Path
from typing import Optional
from .base import Plugin


class MacOSSystemPlugin(Plugin):
    """Plugin to control macOS system-wide dark mode."""

    @property
    def name(self) -> str:
        return "macos-system"

    def validate_config(self) -> tuple[bool, Optional[str]]:
        """Validate that we're running on macOS and wallpaper paths exist."""
        if platform.system() != 'Darwin':
            return False, "macOS system plugin only works on macOS"

        # Check wallpaper paths if provided
        light_wallpaper = self.config.get('light_wallpaper')
        dark_wallpaper = self.config.get('dark_wallpaper')

        if light_wallpaper:
            light_path = Path(light_wallpaper).expanduser()
            if not light_path.exists():
                return False, f"Light wallpaper not found: {light_wallpaper}"

        if dark_wallpaper:
            dark_path = Path(dark_wallpaper).expanduser()
            if not dark_path.exists():
                return False, f"Dark wallpaper not found: {dark_wallpaper}"

        return True, None

    def _set_appearance(self, dark_mode: bool) -> bool:
        """
        Set macOS system appearance using AppleScript.

        Args:
            dark_mode: True for dark mode, False for light mode

        Returns:
            True if successful, False otherwise
        """
        mode_value = "true" if dark_mode else "false"
        applescript = f'''
        tell application "System Events"
            tell appearance preferences
                set dark mode to {mode_value}
            end tell
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

    def _set_wallpaper(self, wallpaper_path: str) -> bool:
        """
        Set macOS desktop wallpaper using AppleScript.

        Args:
            wallpaper_path: Path to the wallpaper image file

        Returns:
            True if successful, False otherwise
        """
        # Expand user path and convert to absolute path
        full_path = Path(wallpaper_path).expanduser().resolve()

        applescript = f'''
        tell application "Finder"
            set desktop picture to POSIX file "{full_path}"
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
                print(f"    Warning (wallpaper): {result.stderr.strip()}")
                return False
            return True
        except subprocess.CalledProcessError as e:
            error_msg = e.stderr.strip() if e.stderr else str(e)
            print(f"    Error (wallpaper): {error_msg}")
            return False
        except subprocess.TimeoutExpired:
            print(f"    Error (wallpaper): Command timed out")
            return False
        except Exception as e:
            print(f"    Error (wallpaper): {e}")
            return False

    def set_light_mode(self) -> bool:
        """Set macOS to light mode and optionally change wallpaper."""
        success = self._set_appearance(dark_mode=False)

        # Change wallpaper if configured
        light_wallpaper = self.config.get('light_wallpaper')
        if light_wallpaper:
            wallpaper_success = self._set_wallpaper(light_wallpaper)
            # Only fail if both failed
            success = success and wallpaper_success

        return success

    def set_dark_mode(self) -> bool:
        """Set macOS to dark mode and optionally change wallpaper."""
        success = self._set_appearance(dark_mode=True)

        # Change wallpaper if configured
        dark_wallpaper = self.config.get('dark_wallpaper')
        if dark_wallpaper:
            wallpaper_success = self._set_wallpaper(dark_wallpaper)
            # Only fail if both failed
            success = success and wallpaper_success

        return success
