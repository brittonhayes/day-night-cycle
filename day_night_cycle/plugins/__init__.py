"""Plugin system for day/night cycle automation."""

from .base import Plugin
from typing import Dict, List, Type, Any
import importlib
import inspect


class PluginManager:
    """Manages loading and executing plugins."""

    def __init__(self):
        self.plugins: List[Plugin] = []
        self._plugin_classes: Dict[str, Type[Plugin]] = {}

    def discover_plugins(self):
        """Automatically discover all available plugins."""
        from . import iterm2, claude_code, cursor

        modules = [iterm2, claude_code, cursor]

        for module in modules:
            for name, obj in inspect.getmembers(module, inspect.isclass):
                if issubclass(obj, Plugin) and obj is not Plugin:
                    plugin_instance = obj({})
                    self._plugin_classes[plugin_instance.name] = obj

    def load_plugins(self, plugin_configs: List[Dict[str, Any]]) -> None:
        """
        Load plugins from configuration.

        Args:
            plugin_configs: List of plugin configuration dictionaries
        """
        self.discover_plugins()

        for config in plugin_configs:
            plugin_name = config.get('name')
            if not plugin_name:
                continue

            plugin_class = self._plugin_classes.get(plugin_name)
            if not plugin_class:
                print(f"Warning: Plugin '{plugin_name}' not found")
                continue

            try:
                plugin = plugin_class(config)
                is_valid, error = plugin.validate_config()

                if not is_valid:
                    print(f"Warning: Plugin '{plugin_name}' config invalid: {error}")
                    continue

                if plugin.is_enabled():
                    self.plugins.append(plugin)
                    print(f"Loaded plugin: {plugin_name}")
            except Exception as e:
                print(f"Error loading plugin '{plugin_name}': {e}")

    def set_light_mode(self) -> Dict[str, bool]:
        """
        Set all enabled plugins to light mode.

        Returns:
            Dictionary mapping plugin names to success status
        """
        results = {}
        for plugin in self.plugins:
            try:
                success = plugin.set_light_mode()
                results[plugin.name] = success
                status = "✓" if success else "✗"
                print(f"  {status} {plugin.name}")
            except Exception as e:
                results[plugin.name] = False
                print(f"  ✗ {plugin.name}: {e}")
        return results

    def set_dark_mode(self) -> Dict[str, bool]:
        """
        Set all enabled plugins to dark mode.

        Returns:
            Dictionary mapping plugin names to success status
        """
        results = {}
        for plugin in self.plugins:
            try:
                success = plugin.set_dark_mode()
                results[plugin.name] = success
                status = "✓" if success else "✗"
                print(f"  {status} {plugin.name}")
            except Exception as e:
                results[plugin.name] = False
                print(f"  ✗ {plugin.name}: {e}")
        return results


__all__ = ['Plugin', 'PluginManager']
