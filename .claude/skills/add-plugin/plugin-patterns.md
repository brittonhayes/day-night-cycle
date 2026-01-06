# Plugin Implementation Patterns

This guide provides detailed patterns for implementing plugins based on how the target application handles theme configuration.

## Pattern 1: JSON Configuration Files

**Use when:** Application stores settings in JSON format (e.g., VS Code, Cursor, Claude Code)

**Example applications:** VS Code, Cursor, Sublime Text, many Electron apps

### Template

```python
"""[App Name] plugin for day/night cycle automation."""

import json
import os
from pathlib import Path
from typing import Optional
from .base import Plugin


class [AppName]Plugin(Plugin):
    """
    Plugin to control [App Name] theme.

    [Add any special notes about how the app handles theme changes]
    """

    def __init__(self, config):
        super().__init__(config)
        # Determine settings path based on OS
        self.settings_path = Path.home() / 'Library' / 'Application Support' / '[App Name]' / 'User' / 'settings.json'

    @property
    def name(self) -> str:
        return "[app-name]"

    def validate_config(self) -> tuple[bool, Optional[str]]:
        """Validate settings file exists."""
        if not self.settings_path.exists():
            return False, f"[App Name] settings not found at {self.settings_path}"
        return True, None

    def _update_theme(self, theme: str) -> bool:
        """
        Update theme in settings.

        Args:
            theme: 'light' or 'dark'

        Returns:
            True if successful, False otherwise
        """
        # Get theme names from config or use defaults
        theme_map = {
            'light': self.config.get('light_theme', '[Default Light Theme]'),
            'dark': self.config.get('dark_theme', '[Default Dark Theme]')
        }

        try:
            # Read current settings
            if self.settings_path.exists():
                with open(self.settings_path, 'r') as f:
                    settings = json.load(f)
            else:
                settings = {}

            target_theme = theme_map[theme]
            current_theme = settings.get('[theme_property_name]')

            # Skip if already set to avoid unnecessary writes
            if current_theme == target_theme:
                return True

            # Update theme property
            settings['[theme_property_name]'] = target_theme

            # Ensure directory exists
            self.settings_path.parent.mkdir(parents=True, exist_ok=True)

            # Write back with explicit flush
            with open(self.settings_path, 'w') as f:
                json.dump(settings, f, indent=2)
                f.flush()
                os.fsync(f.fileno())

            # Verify the write
            with open(self.settings_path, 'r') as f:
                verify_settings = json.load(f)
                if verify_settings.get('[theme_property_name]') != target_theme:
                    print(f"    Warning: Theme not properly saved")
                    return False

            return True
        except FileNotFoundError:
            print(f"    Error: Settings file not found at {self.settings_path}")
            return False
        except json.JSONDecodeError as e:
            print(f"    Error: Invalid JSON in settings file: {e}")
            return False
        except Exception as e:
            print(f"    Error: {e}")
            return False

    def set_light_mode(self) -> bool:
        """Set [App Name] to light theme."""
        return self._update_theme('light')

    def set_dark_mode(self) -> bool:
        """Set [App Name] to dark theme."""
        return self._update_theme('dark')
```

### Key considerations:
- Check if file exists before reading
- Handle malformed JSON gracefully
- Use `os.fsync()` to ensure write completes
- Verify the write succeeded
- Skip writes if already in target mode

## Pattern 2: YAML Configuration Files

**Use when:** Application uses YAML for settings

**Example applications:** Many CLI tools, development tools

### Template

```python
"""[App Name] plugin for day/night cycle automation."""

import yaml
from pathlib import Path
from typing import Optional
from .base import Plugin


class [AppName]Plugin(Plugin):
    """Plugin to control [App Name] theme."""

    def __init__(self, config):
        super().__init__(config)
        self.config_path = Path.home() / '.[app-name]' / 'config.yml'

    @property
    def name(self) -> str:
        return "[app-name]"

    def validate_config(self) -> tuple[bool, Optional[str]]:
        """Validate config file exists."""
        if not self.config_path.exists():
            return False, f"[App Name] config not found at {self.config_path}"
        return True, None

    def _update_theme(self, theme: str) -> bool:
        """Update theme in config."""
        theme_map = {
            'light': self.config.get('light_theme', '[default-light]'),
            'dark': self.config.get('dark_theme', '[default-dark]')
        }

        try:
            # Read current config
            with open(self.config_path, 'r') as f:
                config = yaml.safe_load(f) or {}

            # Update theme
            target_theme = theme_map[theme]
            if config.get('theme') == target_theme:
                return True

            config['theme'] = target_theme

            # Write back
            with open(self.config_path, 'w') as f:
                yaml.dump(config, f, default_flow_style=False)

            return True
        except FileNotFoundError:
            print(f"    Error: Config file not found at {self.config_path}")
            return False
        except yaml.YAMLError as e:
            print(f"    Error: Invalid YAML: {e}")
            return False
        except Exception as e:
            print(f"    Error: {e}")
            return False

    def set_light_mode(self) -> bool:
        """Set [App Name] to light theme."""
        return self._update_theme('light')

    def set_dark_mode(self) -> bool:
        """Set [App Name] to dark theme."""
        return self._update_theme('dark')
```

## Pattern 3: AppleScript Control (macOS)

**Use when:** macOS application supports AppleScript automation

**Example applications:** iTerm2, Terminal, some native macOS apps

### Template

```python
"""[App Name] plugin for day/night cycle automation."""

import subprocess
from typing import Optional
from .base import Plugin


class [AppName]Plugin(Plugin):
    """Plugin to control [App Name] via AppleScript."""

    @property
    def name(self) -> str:
        return "[app-name]"

    def _run_applescript(self, script: str) -> bool:
        """
        Execute AppleScript and return success status.

        Args:
            script: AppleScript code to execute

        Returns:
            True if successful, False otherwise
        """
        try:
            result = subprocess.run(
                ['osascript', '-e', script],
                check=True,
                capture_output=True,
                text=True,
                timeout=5
            )
            return True
        except subprocess.TimeoutExpired:
            print(f"    Error: AppleScript execution timed out")
            return False
        except subprocess.CalledProcessError as e:
            print(f"    Error: AppleScript failed: {e.stderr}")
            return False
        except FileNotFoundError:
            print(f"    Error: osascript not found (not on macOS?)")
            return False
        except Exception as e:
            print(f"    Error: {e}")
            return False

    def set_light_mode(self) -> bool:
        """Set [App Name] to light theme."""
        light_preset = self.config.get('light_preset', '[Default Light]')
        script = f'''
        tell application "[App Name]"
            -- AppleScript commands to set light theme
            -- Example: set current theme to "{light_preset}"
        end tell
        '''
        return self._run_applescript(script)

    def set_dark_mode(self) -> bool:
        """Set [App Name] to dark theme."""
        dark_preset = self.config.get('dark_preset', '[Default Dark]')
        script = f'''
        tell application "[App Name]"
            -- AppleScript commands to set dark theme
            -- Example: set current theme to "{dark_preset}"
        end tell
        '''
        return self._run_applescript(script)
```

### Key considerations:
- Set timeout to prevent hanging
- Capture stderr for error messages
- Handle case where osascript isn't available
- Check if application is running first (optional)

## Pattern 4: Command-Line Interface

**Use when:** Application provides CLI for theme control

**Example applications:** Many terminal apps, CLI-first tools

### Template

```python
"""[App Name] plugin for day/night cycle automation."""

import subprocess
from typing import Optional
from .base import Plugin


class [AppName]Plugin(Plugin):
    """Plugin to control [App Name] via CLI."""

    @property
    def name(self) -> str:
        return "[app-name]"

    def validate_config(self) -> tuple[bool, Optional[str]]:
        """Validate that the CLI tool is available."""
        try:
            result = subprocess.run(
                ['[command-name]', '--version'],
                capture_output=True,
                timeout=5
            )
            if result.returncode != 0:
                return False, "[App Name] CLI not found or not working"
            return True, None
        except FileNotFoundError:
            return False, "[App Name] CLI not found in PATH"
        except Exception as e:
            return False, f"Error checking [App Name] CLI: {e}"

    def _run_command(self, *args) -> bool:
        """Execute CLI command."""
        try:
            result = subprocess.run(
                ['[command-name]', *args],
                check=True,
                capture_output=True,
                text=True,
                timeout=10
            )
            return True
        except subprocess.TimeoutExpired:
            print(f"    Error: Command timed out")
            return False
        except subprocess.CalledProcessError as e:
            print(f"    Error: Command failed: {e.stderr}")
            return False
        except FileNotFoundError:
            print(f"    Error: [command-name] not found in PATH")
            return False
        except Exception as e:
            print(f"    Error: {e}")
            return False

    def set_light_mode(self) -> bool:
        """Set [App Name] to light theme."""
        theme = self.config.get('light_theme', '[default-light]')
        return self._run_command('theme', 'set', theme)

    def set_dark_mode(self) -> bool:
        """Set [App Name] to dark theme."""
        theme = self.config.get('dark_theme', '[default-dark]')
        return self._run_command('theme', 'set', theme)
```

## Pattern 5: Property List (plist) Files (macOS)

**Use when:** macOS application stores preferences in plist format

**Example applications:** Many native macOS applications

### Template

```python
"""[App Name] plugin for day/night cycle automation."""

import subprocess
import plistlib
from pathlib import Path
from typing import Optional
from .base import Plugin


class [AppName]Plugin(Plugin):
    """Plugin to control [App Name] via plist preferences."""

    def __init__(self, config):
        super().__init__(config)
        self.plist_path = Path.home() / 'Library' / 'Preferences' / 'com.[vendor].[app].plist'
        self.defaults_domain = 'com.[vendor].[app]'

    @property
    def name(self) -> str:
        return "[app-name]"

    def validate_config(self) -> tuple[bool, Optional[str]]:
        """Validate plist file exists."""
        if not self.plist_path.exists():
            return False, f"Preferences not found at {self.plist_path}"
        return True, None

    def _update_via_defaults(self, key: str, value: str) -> bool:
        """Update preference using defaults command."""
        try:
            subprocess.run(
                ['defaults', 'write', self.defaults_domain, key, value],
                check=True,
                capture_output=True,
                timeout=5
            )
            return True
        except Exception as e:
            print(f"    Error: {e}")
            return False

    def _update_theme(self, theme: str) -> bool:
        """Update theme preference."""
        theme_map = {
            'light': self.config.get('light_theme', '[LightTheme]'),
            'dark': self.config.get('dark_theme', '[DarkTheme]')
        }

        theme_value = theme_map[theme]
        success = self._update_via_defaults('[ThemeKey]', theme_value)

        if success:
            # Optional: Notify app to reload preferences
            # Some apps need this, others watch for changes automatically
            pass

        return success

    def set_light_mode(self) -> bool:
        """Set [App Name] to light theme."""
        return self._update_theme('light')

    def set_dark_mode(self) -> bool:
        """Set [App Name] to dark theme."""
        return self._update_theme('dark')
```

## Pattern 6: Multiple Configuration Methods

**Use when:** Application can be controlled via multiple methods (prefer most reliable)

### Decision priority:
1. **Native API** (if available and documented)
2. **Configuration file** (most reliable, works when app isn't running)
3. **AppleScript** (good for macOS, requires app to be running)
4. **CLI** (if official CLI exists)
5. **Environment variables** (for terminal apps)

## Common Helper Methods

### Check if Application is Running

```python
def _is_app_running(self, app_name: str) -> bool:
    """Check if application is currently running (macOS)."""
    try:
        result = subprocess.run(
            ['pgrep', '-x', app_name],
            capture_output=True
        )
        return result.returncode == 0
    except:
        return False
```

### Get Available Themes

```python
def get_available_themes(self) -> list[str]:
    """
    Get list of available themes.
    Useful for creating helper scripts.
    """
    try:
        # Implementation depends on how app exposes themes
        pass
    except Exception as e:
        print(f"Error: {e}")
        return []
```

## Research Checklist

When researching a new application, find:

- [ ] Settings file location and format
- [ ] Theme property name(s)
- [ ] Available theme names/identifiers
- [ ] Does app auto-reload settings?
- [ ] Are there platform differences? (macOS/Linux/Windows)
- [ ] Any special permissions needed?
- [ ] Official documentation on automation
- [ ] Existing automation examples in community

## Testing Checklist

After implementation:

- [ ] Plugin imports without errors
- [ ] `validate_config()` works correctly
- [ ] Light mode switches successfully
- [ ] Dark mode switches successfully
- [ ] Error messages are helpful
- [ ] Returns correct boolean values
- [ ] Handles missing files gracefully
- [ ] Works when app is not running (if applicable)
- [ ] Works when app is running (if applicable)
