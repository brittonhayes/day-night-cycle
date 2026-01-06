#!/usr/bin/env python3
"""Update launchd schedule based on sunrise/sunset times."""

import sys
from pathlib import Path
from datetime import datetime
import plistlib

# Add parent directory to path
sys.path.insert(0, str(Path(__file__).parent.parent))

from day_night_cycle.scheduler import SolarScheduler
from day_night_cycle.main import DayNightCycle


def create_calendar_interval(dt: datetime) -> dict:
    """Create a launchd calendar interval from datetime."""
    return {
        'Month': dt.month,
        'Day': dt.day,
        'Hour': dt.hour,
        'Minute': dt.minute,
    }


def update_launchd_schedule():
    """Update launchd plists with today's sunrise/sunset times."""
    # Load app to get scheduler
    app = DayNightCycle()
    solar_times = app.scheduler.get_solar_times()

    sunrise = solar_times['sunrise']
    sunset = solar_times['sunset']

    # Path to launchd plist files
    launchd_dir = Path.home() / 'Library' / 'LaunchAgents'
    sunrise_plist = launchd_dir / 'com.daynightcycle.sunrise.plist'
    sunset_plist = launchd_dir / 'com.daynightcycle.sunset.plist'

    project_dir = Path(__file__).parent.parent.absolute()
    python_path = sys.executable

    # Create sunrise plist
    sunrise_config = {
        'Label': 'com.daynightcycle.sunrise',
        'ProgramArguments': [
            python_path,
            '-m',
            'day_night_cycle',
            'light',
        ],
        'WorkingDirectory': str(project_dir),
        'StartCalendarInterval': create_calendar_interval(sunrise),
        'StandardOutPath': str(project_dir / 'logs' / 'sunrise.log'),
        'StandardErrorPath': str(project_dir / 'logs' / 'sunrise.error.log'),
        'RunAtLoad': False,
    }

    # Create sunset plist
    sunset_config = {
        'Label': 'com.daynightcycle.sunset',
        'ProgramArguments': [
            python_path,
            '-m',
            'day_night_cycle',
            'dark',
        ],
        'WorkingDirectory': str(project_dir),
        'StartCalendarInterval': create_calendar_interval(sunset),
        'StandardOutPath': str(project_dir / 'logs' / 'sunset.log'),
        'StandardErrorPath': str(project_dir / 'logs' / 'sunset.error.log'),
        'RunAtLoad': False,
    }

    # Ensure launchd directory exists
    launchd_dir.mkdir(parents=True, exist_ok=True)

    # Write plists
    with open(sunrise_plist, 'wb') as f:
        plistlib.dump(sunrise_config, f)

    with open(sunset_plist, 'wb') as f:
        plistlib.dump(sunset_config, f)

    print(f"Updated schedule:")
    print(f"  Sunrise: {app.scheduler.format_time(sunrise)}")
    print(f"  Sunset: {app.scheduler.format_time(sunset)}")
    print(f"\nPlist files created at:")
    print(f"  {sunrise_plist}")
    print(f"  {sunset_plist}")


if __name__ == '__main__':
    update_launchd_schedule()
