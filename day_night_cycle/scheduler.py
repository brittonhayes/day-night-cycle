"""Sunrise/sunset calculation and scheduling logic."""

from datetime import datetime, timedelta
from typing import Tuple, Optional
from zoneinfo import ZoneInfo
from astral import LocationInfo
from astral.sun import sun


class SolarScheduler:
    """Calculates sunrise and sunset times based on location."""

    def __init__(self, latitude: float, longitude: float, timezone: str, location_name: str = "Custom"):
        """
        Initialize the solar scheduler.

        Args:
            latitude: Latitude in decimal degrees
            longitude: Longitude in decimal degrees
            timezone: Timezone string (e.g., 'America/Los_Angeles')
            location_name: Optional name for the location
        """
        self.location = LocationInfo(
            name=location_name,
            region="",
            timezone=timezone,
            latitude=latitude,
            longitude=longitude
        )
        self.timezone = ZoneInfo(timezone)

    def get_solar_times(self, date: Optional[datetime] = None) -> dict:
        """
        Get solar times for a specific date.

        Args:
            date: Date to calculate for (defaults to today)

        Returns:
            Dictionary with sunrise, sunset, and other solar times
        """
        if date is None:
            date = datetime.now(self.timezone)

        solar_times = sun(self.location.observer, date=date.date(), tzinfo=self.timezone)

        return {
            'sunrise': solar_times['sunrise'],
            'sunset': solar_times['sunset'],
            'dawn': solar_times['dawn'],
            'dusk': solar_times['dusk'],
            'noon': solar_times['noon']
        }

    def get_next_transition(self) -> Tuple[datetime, str]:
        """
        Get the next solar transition (sunrise or sunset).

        Returns:
            Tuple of (next_transition_time, transition_type)
            where transition_type is 'sunrise' or 'sunset'
        """
        now = datetime.now(self.timezone)
        solar_times = self.get_solar_times(now)

        sunrise = solar_times['sunrise']
        sunset = solar_times['sunset']

        # If both transitions are in the past, get tomorrow's times
        if now > sunset:
            tomorrow = now + timedelta(days=1)
            solar_times = self.get_solar_times(tomorrow)
            return solar_times['sunrise'], 'sunrise'

        # Return the next upcoming transition
        if now < sunrise:
            return sunrise, 'sunrise'
        else:
            return sunset, 'sunset'

    def is_daytime(self) -> bool:
        """
        Check if it's currently daytime.

        Returns:
            True if between sunrise and sunset, False otherwise
        """
        now = datetime.now(self.timezone)
        solar_times = self.get_solar_times(now)

        sunrise = solar_times['sunrise']
        sunset = solar_times['sunset']

        return sunrise <= now <= sunset

    def format_time(self, dt: datetime) -> str:
        """Format datetime for display."""
        return dt.strftime('%I:%M %p')

    def get_schedule_summary(self) -> str:
        """Get a human-readable schedule summary."""
        solar_times = self.get_solar_times()
        next_transition, transition_type = self.get_next_transition()

        sunrise_str = self.format_time(solar_times['sunrise'])
        sunset_str = self.format_time(solar_times['sunset'])
        next_str = self.format_time(next_transition)

        current_mode = "light" if self.is_daytime() else "dark"

        return f"""
Current mode: {current_mode}
Today's sunrise: {sunrise_str}
Today's sunset: {sunset_str}
Next transition: {next_str} ({transition_type})
"""
