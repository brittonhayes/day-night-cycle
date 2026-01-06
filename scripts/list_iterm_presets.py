#!/usr/bin/env python3
"""List available iTerm2 color presets."""

import plistlib
import os
import sys

def list_presets():
    """List all available iTerm2 color presets."""
    plist_path = os.path.expanduser('~/Library/Preferences/com.googlecode.iterm2.plist')

    if not os.path.exists(plist_path):
        print("iTerm2 preferences not found.")
        print("Make sure iTerm2 is installed and has been run at least once.")
        return 1

    try:
        with open(plist_path, 'rb') as f:
            prefs = plistlib.load(f)

        print("=" * 60)
        print("Available iTerm2 Color Presets")
        print("=" * 60)
        print()

        # List custom presets
        if 'Custom Color Presets' in prefs and prefs['Custom Color Presets']:
            print("Custom Presets:")
            for preset_name in sorted(prefs['Custom Color Presets'].keys()):
                print(f"  • {preset_name}")
            print()

        # List built-in presets (common ones)
        print("Built-in Presets (partial list):")
        builtin_presets = [
            "Dark Background",
            "Light Background",
            "Pastel (Dark Background)",
            "Solarized Dark",
            "Solarized Light",
            "Tango Dark",
            "Tango Light"
        ]
        for preset in builtin_presets:
            print(f"  • {preset}")

        print()
        print("To use a preset, add it to your config.yaml:")
        print("  light_preset: \"Light Background\"")
        print("  dark_preset: \"githubdark\"")
        print()

        return 0

    except Exception as e:
        print(f"Error reading iTerm2 preferences: {e}")
        return 1

if __name__ == '__main__':
    sys.exit(list_presets())
