# day-night-cycle

Automatically switch your apps between light and dark themes based on sunrise and sunset for your location.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/brittonhayes/day-night-cycle/main/install.sh | bash
```

## Supported Apps

- **iTerm2** - Switches color presets
- **Claude Code** - Switches theme settings
- **Cursor** - Switches color themes
- **Neovim** - Switches background setting
- **macOS System** - Switches system appearance

## Usage

After installation, themes switch automatically at sunrise and sunset. Manual control:

```bash
# Apply mode based on current time
python3 -m day_night_cycle auto

# Force light or dark mode
python3 -m day_night_cycle light
python3 -m day_night_cycle dark

# Show status
python3 -m day_night_cycle status
```

## Configuration

Edit `~/.config/day-night-cycle/config.yaml`:

```yaml
location:
  latitude: 46.0645    # Your coordinates
  longitude: -118.3430
  timezone: "America/Los_Angeles"

plugins:
  - name: iterm2
    enabled: true
    light_preset: "Light Background"
    dark_preset: "Dark Background"

  - name: cursor
    enabled: true
    light_theme: "GitHub Light"
    dark_theme: "Dark Modern"
```

Find coordinates at [latlong.net](https://www.latlong.net/) and timezones at [Wikipedia](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).

## Custom Plugins

Create `day_night_cycle/plugins/myapp.py`:

```python
from .base import Plugin

class MyAppPlugin(Plugin):
    @property
    def name(self) -> str:
        return "myapp"

    def set_light_mode(self) -> bool:
        # Switch to light mode
        return True

    def set_dark_mode(self) -> bool:
        # Switch to dark mode
        return True
```

Import in `day_night_cycle/plugins/__init__.py` and add to config.

## Uninstall

```bash
~/.config/day-night-cycle/scripts/uninstall.sh
```

## License

This project is released into the public domain under the [Unlicense](.github/LICENSE).
