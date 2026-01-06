#!/usr/bin/env python3
"""Test configuration and show what would happen."""

import sys
from pathlib import Path

# Add parent directory to path
sys.path.insert(0, str(Path(__file__).parent.parent))

from day_night_cycle.main import DayNightCycle


def main():
    """Test the configuration and display status."""
    print("=" * 60)
    print("Day/Night Cycle - Configuration Test")
    print("=" * 60)
    print()

    try:
        app = DayNightCycle()
    except SystemExit:
        print("\nConfiguration test failed.")
        print("Please ensure config.yaml is properly configured.")
        return 1

    # Show schedule
    print(app.scheduler.get_schedule_summary())

    # Show plugins
    print("\nConfigured plugins:")
    if not app.plugin_manager.plugins:
        print("  (none)")
    else:
        for plugin in app.plugin_manager.plugins:
            print(f"  ✓ {plugin.name}")

    # Check what mode should be active
    is_day = app.scheduler.is_daytime()
    current_mode = "light" if is_day else "dark"

    print(f"\nCurrent mode should be: {current_mode}")
    print("\nConfiguration test passed!")
    print("\nRun './scripts/install.sh' to install the automation.")
    print()

    return 0


if __name__ == '__main__':
    sys.exit(main())
