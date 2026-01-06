"""Base plugin interface for day/night cycle automation."""

from abc import ABC, abstractmethod
from typing import Dict, Any, Optional


class Plugin(ABC):
    """Base plugin interface that all plugins must implement."""

    def __init__(self, config: Dict[str, Any]):
        """
        Initialize the plugin with configuration.

        Args:
            config: Plugin-specific configuration dictionary
        """
        self.config = config
        self.enabled = config.get('enabled', True)

    @property
    @abstractmethod
    def name(self) -> str:
        """Return the unique name of the plugin."""
        pass

    @abstractmethod
    def set_light_mode(self) -> bool:
        """
        Set the application to light mode.

        Returns:
            True if successful, False otherwise
        """
        pass

    @abstractmethod
    def set_dark_mode(self) -> bool:
        """
        Set the application to dark mode.

        Returns:
            True if successful, False otherwise
        """
        pass

    def is_enabled(self) -> bool:
        """Check if the plugin is enabled."""
        return self.enabled

    def validate_config(self) -> tuple[bool, Optional[str]]:
        """
        Validate the plugin configuration.

        Returns:
            Tuple of (is_valid, error_message)
        """
        return True, None
