"""Claude Code plugin for day/night cycle automation."""

import json
import os
from pathlib import Path
from typing import Optional
from .base import Plugin


class ClaudeCodePlugin(Plugin):
    """Plugin to control Claude Code theme."""

    def __init__(self, config):
        super().__init__(config)
        # Claude Code settings path
        self.settings_path = Path.home() / '.claude' / 'settings.json'

    @property
    def name(self) -> str:
        return "claude-code"

    def validate_config(self) -> tuple[bool, Optional[str]]:
        """Validate Claude Code settings file exists."""
        if not self.settings_path.exists():
            return False, f"Claude Code settings not found at {self.settings_path}"
        return True, None

    def _update_theme(self, theme: str) -> bool:
        """
        Update Claude Code theme in settings.

        Args:
            theme: Theme name ('light' or 'dark')

        Returns:
            True if successful, False otherwise
        """
        try:
            # Read current settings
            with open(self.settings_path, 'r') as f:
                settings = json.load(f)

            # Update theme
            settings['theme'] = theme

            # Write back
            with open(self.settings_path, 'w') as f:
                json.dump(settings, f, indent=2)

            return True
        except Exception as e:
            print(f"    Error: {e}")
            return False

    def set_light_mode(self) -> bool:
        """Set Claude Code to light theme."""
        return self._update_theme('light')

    def set_dark_mode(self) -> bool:
        """Set Claude Code to dark theme."""
        return self._update_theme('dark')
