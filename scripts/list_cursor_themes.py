#!/usr/bin/env python3
"""List available Cursor/VS Code themes."""

import json
import os
from pathlib import Path
import sys


def list_themes():
    """List installed themes in Cursor."""
    # Cursor settings path
    settings_path = Path.home() / 'Library' / 'Application Support' / 'Cursor' / 'User' / 'settings.json'
    extensions_path = Path.home() / 'Library' / 'Application Support' / 'Cursor' / 'User' / 'extensions'

    print("=" * 60)
    print("Cursor Theme Configuration")
    print("=" * 60)
    print()

    # Show current theme
    if settings_path.exists():
        try:
            with open(settings_path, 'r') as f:
                settings = json.load(f)
            current_theme = settings.get('workbench.colorTheme', 'Not set')
            print(f"Current theme: {current_theme}")
            print()
        except Exception as e:
            print(f"Error reading settings: {e}")
            print()

    # List common built-in themes
    print("Common Built-in Themes:")
    builtin_themes = [
        "Dark+ (default dark)",
        "Dark Modern",
        "Dark High Contrast",
        "Light+ (default light)",
        "Light Modern",
        "Light High Contrast",
        "Solarized Dark",
        "Solarized Light",
        "Monokai",
        "GitHub Light",
        "GitHub Dark",
        "GitHub Dark Dimmed",
    ]
    for theme in builtin_themes:
        print(f"  • {theme}")

    print()
    print("To see all installed themes:")
    print("  1. Open Cursor")
    print("  2. Press Cmd+K Cmd+T (or Cmd+Shift+P, then 'Color Theme')")
    print("  3. Browse the list and note the exact name")
    print()
    print("Then update config.yaml with the exact theme names:")
    print('  light_theme: "GitHub Light"')
    print('  dark_theme: "Dark Modern"')
    print()

    return 0


if __name__ == '__main__':
    sys.exit(list_themes())
