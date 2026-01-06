"""Main entry point for day/night cycle automation."""

import argparse
import sys
from pathlib import Path
from typing import Optional
import yaml

from .scheduler import SolarScheduler
from .plugins import PluginManager


class DayNightCycle:
    """Main application for day/night cycle automation."""

    def __init__(self, config_path: Optional[Path] = None):
        """
        Initialize the application.

        Args:
            config_path: Path to configuration file
        """
        if config_path is None:
            config_path = Path(__file__).parent.parent / 'config.yaml'

        self.config_path = config_path
        self.config = self._load_config()
        self.scheduler = self._init_scheduler()
        self.plugin_manager = self._init_plugins()

    def _load_config(self) -> dict:
        """Load configuration from file."""
        if not self.config_path.exists():
            print(f"Error: Configuration file not found at {self.config_path}")
            sys.exit(1)

        try:
            with open(self.config_path, 'r') as f:
                return yaml.safe_load(f)
        except Exception as e:
            print(f"Error loading configuration: {e}")
            sys.exit(1)

    def _init_scheduler(self) -> SolarScheduler:
        """Initialize the solar scheduler from config."""
        location = self.config.get('location', {})

        latitude = location.get('latitude')
        longitude = location.get('longitude')
        timezone = location.get('timezone')

        if not all([latitude, longitude, timezone]):
            print("Error: location (latitude, longitude, timezone) must be configured")
            sys.exit(1)

        return SolarScheduler(
            latitude=latitude,
            longitude=longitude,
            timezone=timezone,
            location_name=location.get('name', 'Custom')
        )

    def _init_plugins(self) -> PluginManager:
        """Initialize the plugin manager."""
        manager = PluginManager()
        plugin_configs = self.config.get('plugins', [])
        manager.load_plugins(plugin_configs)
        return manager

    def apply_mode(self, mode: str) -> None:
        """
        Apply light or dark mode to all plugins.

        Args:
            mode: 'light' or 'dark'
        """
        print(f"\nApplying {mode} mode...")

        if mode == 'light':
            results = self.plugin_manager.set_light_mode()
        elif mode == 'dark':
            results = self.plugin_manager.set_dark_mode()
        else:
            print(f"Error: Unknown mode '{mode}'")
            return

        success_count = sum(1 for success in results.values() if success)
        total_count = len(results)

        print(f"\nCompleted: {success_count}/{total_count} plugins successful")

    def auto_apply(self) -> None:
        """Automatically apply mode based on current time."""
        if self.scheduler.is_daytime():
            self.apply_mode('light')
        else:
            self.apply_mode('dark')

    def show_status(self) -> None:
        """Show current status and schedule."""
        print(self.scheduler.get_schedule_summary())

        print("Configured plugins:")
        for plugin in self.plugin_manager.plugins:
            print(f"  • {plugin.name}")

    def show_next(self) -> None:
        """Show next scheduled transition."""
        next_time, transition_type = self.scheduler.get_next_transition()
        print(f"Next transition: {self.scheduler.format_time(next_time)} ({transition_type})")


def main():
    """CLI entry point."""
    parser = argparse.ArgumentParser(
        description='Day/Night Cycle Automation - Automatically switch themes based on sunrise/sunset'
    )
    parser.add_argument(
        '--config',
        type=Path,
        help='Path to configuration file (default: config.yaml)'
    )
    parser.add_argument(
        'command',
        choices=['auto', 'light', 'dark', 'status', 'next'],
        help='Command to execute'
    )

    args = parser.parse_args()

    app = DayNightCycle(config_path=args.config)

    if args.command == 'auto':
        app.auto_apply()
    elif args.command == 'light':
        app.apply_mode('light')
    elif args.command == 'dark':
        app.apply_mode('dark')
    elif args.command == 'status':
        app.show_status()
    elif args.command == 'next':
        app.show_next()


if __name__ == '__main__':
    main()
