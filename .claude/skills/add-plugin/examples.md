# Example Plugin Implementations

This document shows complete examples of adding plugins for common applications.

## Example 1: Visual Studio Code

**User request:** "Add a plugin for VS Code"

### Research findings:
- Settings location: `~/Library/Application Support/Code/User/settings.json`
- Theme property: `workbench.colorTheme`
- Popular themes: "Light+ (default light)", "Dark+ (default dark)"
- Auto-reload: Yes, VS Code watches settings.json

### Implementation:

**File:** `day_night_cycle/plugins/vscode.py`

```python
"""Visual Studio Code plugin for day/night cycle automation."""

import json
import os
from pathlib import Path
from typing import Optional
from .base import Plugin


class VSCodePlugin(Plugin):
    """
    Plugin to control VS Code theme.

    VS Code watches its settings.json and applies changes automatically.
    """

    def __init__(self, config):
        super().__init__(config)
        self.settings_path = Path.home() / 'Library' / 'Application Support' / 'Code' / 'User' / 'settings.json'

    @property
    def name(self) -> str:
        return "vscode"

    def validate_config(self) -> tuple[bool, Optional[str]]:
        """Validate VS Code settings file exists."""
        if not self.settings_path.exists():
            return False, f"VS Code settings not found at {self.settings_path}"
        return True, None

    def _update_theme(self, theme: str) -> bool:
        """Update VS Code theme in settings."""
        theme_map = {
            'light': self.config.get('light_theme', 'Light+ (default light)'),
            'dark': self.config.get('dark_theme', 'Dark+ (default dark)')
        }

        try:
            if self.settings_path.exists():
                with open(self.settings_path, 'r') as f:
                    settings = json.load(f)
            else:
                settings = {}

            target_theme = theme_map[theme]
            current_theme = settings.get('workbench.colorTheme')

            if current_theme == target_theme:
                return True

            settings['workbench.colorTheme'] = target_theme

            self.settings_path.parent.mkdir(parents=True, exist_ok=True)

            with open(self.settings_path, 'w') as f:
                json.dump(settings, f, indent=2)
                f.flush()
                os.fsync(f.fileno())

            with open(self.settings_path, 'r') as f:
                verify_settings = json.load(f)
                if verify_settings.get('workbench.colorTheme') != target_theme:
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
        """Set VS Code to light theme."""
        return self._update_theme('light')

    def set_dark_mode(self) -> bool:
        """Set VS Code to dark theme."""
        return self._update_theme('dark')
```

### Integration:

**Update:** `day_night_cycle/plugins/__init__.py`
```python
from . import vscode
```

### Configuration:

**Add to:** `config.yaml`
```yaml
plugins:
  - name: vscode
    enabled: true
    light_theme: "GitHub Light"
    dark_theme: "GitHub Dark"
```

### Testing:
```bash
python3 -m day_night_cycle light  # Should switch VS Code to light theme
python3 -m day_night_cycle dark   # Should switch VS Code to dark theme
```

## Example 2: Kitty Terminal

**User request:** "Add support for Kitty terminal"

### Research findings:
- Settings location: `~/.config/kitty/kitty.conf`
- Theme approach: Include theme files
- Theme location: `~/.config/kitty/themes/`
- Reload command: `kill -SIGUSR1 $(pgrep kitty)`

### Implementation:

**File:** `day_night_cycle/plugins/kitty.py`

```python
"""Kitty terminal plugin for day/night cycle automation."""

import subprocess
from pathlib import Path
from typing import Optional
from .base import Plugin


class KittyPlugin(Plugin):
    """Plugin to control Kitty terminal theme."""

    def __init__(self, config):
        super().__init__(config)
        self.config_path = Path.home() / '.config' / 'kitty' / 'kitty.conf'
        self.light_theme = self.config.get('light_theme', 'light')
        self.dark_theme = self.config.get('dark_theme', 'dark')

    @property
    def name(self) -> str:
        return "kitty"

    def validate_config(self) -> tuple[bool, Optional[str]]:
        """Validate Kitty config exists."""
        if not self.config_path.exists():
            return False, f"Kitty config not found at {self.config_path}"
        return True, None

    def _update_theme(self, theme_name: str) -> bool:
        """Update Kitty theme by modifying config."""
        try:
            # Read current config
            with open(self.config_path, 'r') as f:
                lines = f.readlines()

            # Remove old theme includes
            new_lines = [
                line for line in lines
                if not line.strip().startswith('include themes/')
            ]

            # Add new theme include
            new_lines.append(f'include themes/{theme_name}.conf\n')

            # Write back
            with open(self.config_path, 'w') as f:
                f.writelines(new_lines)

            # Reload Kitty if running
            self._reload_kitty()

            return True
        except FileNotFoundError:
            print(f"    Error: Config file not found at {self.config_path}")
            return False
        except Exception as e:
            print(f"    Error: {e}")
            return False

    def _reload_kitty(self) -> None:
        """Send reload signal to Kitty if running."""
        try:
            # Find Kitty process
            result = subprocess.run(
                ['pgrep', 'kitty'],
                capture_output=True,
                text=True
            )
            if result.returncode == 0 and result.stdout.strip():
                pid = result.stdout.strip().split()[0]
                # Send SIGUSR1 to reload config
                subprocess.run(['kill', '-SIGUSR1', pid], check=False)
        except:
            # If reload fails, it's okay - user can restart Kitty
            pass

    def set_light_mode(self) -> bool:
        """Set Kitty to light theme."""
        return self._update_theme(self.light_theme)

    def set_dark_mode(self) -> bool:
        """Set Kitty to dark theme."""
        return self._update_theme(self.dark_theme)
```

### Configuration:

```yaml
plugins:
  - name: kitty
    enabled: true
    light_theme: "light"  # Corresponds to themes/light.conf
    dark_theme: "dark"    # Corresponds to themes/dark.conf
```

## Example 3: Slack Desktop

**User request:** "Add Slack to the day/night cycle"

### Research findings:
- Settings location: `~/Library/Application Support/Slack/storage/slack-settings`
- Format: JSON (actually a state file)
- Theme property: `theme` (values: "light", "dark", "system")
- Requires app restart to take effect

### Implementation:

**File:** `day_night_cycle/plugins/slack.py`

```python
"""Slack desktop plugin for day/night cycle automation."""

import json
from pathlib import Path
from typing import Optional
from .base import Plugin


class SlackPlugin(Plugin):
    """
    Plugin to control Slack desktop theme.

    Note: Slack requires a restart to apply theme changes.
    """

    def __init__(self, config):
        super().__init__(config)
        self.settings_path = Path.home() / 'Library' / 'Application Support' / 'Slack' / 'storage' / 'slack-settings'

    @property
    def name(self) -> str:
        return "slack"

    def validate_config(self) -> tuple[bool, Optional[str]]:
        """Validate Slack settings exist."""
        if not self.settings_path.exists():
            return False, f"Slack settings not found at {self.settings_path}"
        return True, None

    def _update_theme(self, theme: str) -> bool:
        """Update Slack theme."""
        try:
            # Read current settings
            with open(self.settings_path, 'r') as f:
                settings = json.load(f)

            # Update theme
            if settings.get('theme') == theme:
                return True

            settings['theme'] = theme

            # Write back
            with open(self.settings_path, 'w') as f:
                json.dump(settings, f, indent=2)

            print(f"    Note: Please restart Slack to apply theme change")
            return True
        except FileNotFoundError:
            print(f"    Error: Settings file not found")
            return False
        except json.JSONDecodeError as e:
            print(f"    Error: Invalid JSON: {e}")
            return False
        except Exception as e:
            print(f"    Error: {e}")
            return False

    def set_light_mode(self) -> bool:
        """Set Slack to light theme."""
        return self._update_theme('light')

    def set_dark_mode(self) -> bool:
        """Set Slack to dark theme."""
        return self._update_theme('dark')
```

### Configuration:

```yaml
plugins:
  - name: slack
    enabled: true
```

### Note:
This plugin requires users to restart Slack after the theme change. Consider adding a helper script to automate the restart if needed.

## Example 4: Vim/Neovim

**User request:** "Add Vim/Neovim support"

### Research findings:
- Settings location: `~/.vimrc` or `~/.config/nvim/init.vim`
- Theme setting: `colorscheme [name]`
- No automatic reload (unless using autocmd)

### Implementation:

**File:** `day_night_cycle/plugins/vim.py`

```python
"""Vim/Neovim plugin for day/night cycle automation."""

from pathlib import Path
from typing import Optional
from .base import Plugin


class VimPlugin(Plugin):
    """Plugin to control Vim/Neovim colorscheme."""

    def __init__(self, config):
        super().__init__(config)
        # Support both Vim and Neovim
        self.vim_rc = Path.home() / '.vimrc'
        self.nvim_rc = Path.home() / '.config' / 'nvim' / 'init.vim'

        # Determine which one to use
        if self.nvim_rc.exists():
            self.config_path = self.nvim_rc
        elif self.vim_rc.exists():
            self.config_path = self.vim_rc
        else:
            self.config_path = self.vim_rc  # Default for creation

    @property
    def name(self) -> str:
        return "vim"

    def _update_colorscheme(self, colorscheme: str) -> bool:
        """Update Vim colorscheme in config file."""
        try:
            # Read current config
            if self.config_path.exists():
                content = self.config_path.read_text()
                lines = content.split('\n')
            else:
                lines = []

            # Remove old colorscheme lines
            new_lines = [
                line for line in lines
                if not line.strip().startswith('colorscheme ')
            ]

            # Add new colorscheme
            new_lines.append(f'colorscheme {colorscheme}')

            # Write back
            self.config_path.parent.mkdir(parents=True, exist_ok=True)
            self.config_path.write_text('\n'.join(new_lines))

            print(f"    Note: Restart Vim/Neovim or run :source % to apply")
            return True
        except Exception as e:
            print(f"    Error: {e}")
            return False

    def set_light_mode(self) -> bool:
        """Set Vim to light colorscheme."""
        colorscheme = self.config.get('light_colorscheme', 'morning')
        return self._update_colorscheme(colorscheme)

    def set_dark_mode(self) -> bool:
        """Set Vim to dark colorscheme."""
        colorscheme = self.config.get('dark_colorscheme', 'evening')
        return self._update_colorscheme(colorscheme)
```

### Configuration:

```yaml
plugins:
  - name: vim
    enabled: true
    light_colorscheme: "solarized"
    dark_colorscheme: "gruvbox"
```

## Helper Script Example

For applications with discoverable themes, create helper scripts:

**File:** `scripts/list_vscode_themes.py`

```python
#!/usr/bin/env python3
"""List available VS Code themes."""

import json
from pathlib import Path

settings_path = Path.home() / 'Library' / 'Application Support' / 'Code' / 'User' / 'settings.json'

if not settings_path.exists():
    print("VS Code settings not found")
    exit(1)

with open(settings_path) as f:
    settings = json.load(f)

current = settings.get('workbench.colorTheme', 'Not set')
print(f"Current theme: {current}\n")
print("To see all available themes:")
print("1. Open VS Code")
print("2. Press Cmd+K Cmd+T (or Ctrl+K Ctrl+T)")
print("3. Browse and note the exact theme name")
print("\nPopular themes:")
print("  - Light: 'Light+ (default light)', 'GitHub Light', 'Solarized Light'")
print("  - Dark: 'Dark+ (default dark)', 'GitHub Dark', 'Monokai'")
```

## Common Patterns Observed

1. **JSON config files** (VS Code, Cursor, Electron apps)
   - Most reliable method
   - Usually auto-reload
   - Need exact theme names

2. **Config file includes** (Kitty, Alacritty)
   - Include external theme files
   - May need manual reload signal

3. **Property-based configs** (.vimrc, .bashrc)
   - Update specific lines
   - Usually need app restart

4. **State files** (Slack, some apps)
   - JSON but not traditional settings
   - Often require app restart

## Tips for Quick Implementation

1. **Start with research** - Don't guess file locations
2. **Check existing plugins** - Follow established patterns
3. **Test error cases** - What if file doesn't exist?
4. **Verify writes** - Did the change actually save?
5. **Note reload requirements** - Does app need restart?
6. **Use exact theme names** - Copy from app's UI
7. **Handle platform differences** - Windows/Linux/macOS paths differ
8. **Provide helpful errors** - Guide users to fix issues
