"""Cursor plugin for day/night cycle automation."""

import json
import os
from pathlib import Path
from typing import Optional
from .base import Plugin


class CursorPlugin(Plugin):
    """
    Plugin to control Cursor theme.

    Cursor watches its settings.json file and applies theme changes
    automatically within 1-2 seconds while running.
    """

    def __init__(self, config):
        super().__init__(config)
        # Cursor settings path (similar to VS Code)
        self.settings_path = Path.home() / 'Library' / 'Application Support' / 'Cursor' / 'User' / 'settings.json'

    @property
    def name(self) -> str:
        return "cursor"

    def validate_config(self) -> tuple[bool, Optional[str]]:
        """Validate Cursor settings file exists."""
        if not self.settings_path.exists():
            return False, f"Cursor settings not found at {self.settings_path}"
        return True, None

    def _update_theme(self, theme: str) -> bool:
        """
        Update Cursor theme in settings.

        Args:
            theme: Theme identifier

        Returns:
            True if successful, False otherwise
        """
        # Theme mappings
        theme_map = {
            'light': self.config.get('light_theme', 'Default Light+'),
            'dark': self.config.get('dark_theme', 'Default Dark+')
        }

        try:
            # Read current settings
            if self.settings_path.exists():
                with open(self.settings_path, 'r') as f:
                    settings = json.load(f)
            else:
                settings = {}

            target_theme = theme_map[theme]
            current_theme = settings.get('workbench.colorTheme')

            # Check if already set to avoid unnecessary writes
            if current_theme == target_theme:
                return True

            # Update theme
            settings['workbench.colorTheme'] = target_theme

            # Ensure directory exists
            self.settings_path.parent.mkdir(parents=True, exist_ok=True)

            # Write back with explicit flush to ensure file is updated
            with open(self.settings_path, 'w') as f:
                json.dump(settings, f, indent=2)
                f.flush()
                os.fsync(f.fileno())

            # Verify the write was successful
            with open(self.settings_path, 'r') as f:
                verify_settings = json.load(f)
                if verify_settings.get('workbench.colorTheme') != target_theme:
                    print(f"    Warning: Theme not properly saved")
                    return False

            return True
        except Exception as e:
            print(f"    Error: {e}")
            return False

    def set_light_mode(self) -> bool:
        """Set Cursor to light theme."""
        return self._update_theme('light')

    def set_dark_mode(self) -> bool:
        """Set Cursor to dark theme."""
        return self._update_theme('dark')
