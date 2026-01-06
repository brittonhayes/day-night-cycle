# Day/Night Cycle Automation

Automatically switch application themes between light and dark modes based on sunrise and sunset times for your location.

## Features

- **Solar-based scheduling**: Privately uses your location to calculate precise sunrise and sunset times
- **Plugin architecture**: Easily extensible system for adding new applications
- **Automatic scheduling**: Uses macOS launchd for reliable background execution
- **Multiple applications**: Built-in support for iTerm2, Claude Code, and Cursor
- **Zero maintenance**: Automatically updates schedule daily

## Built-in Plugins

- **iTerm2**: Switches color presets
- **Claude Code**: Switches theme settings
- **Cursor**: Switches color themes

## Installation

1. **Clone or download this repository**

2. **Configure your location**

   ```bash
   cp config.example.yaml config.yaml
   ```

   Edit `config.yaml` and set your coordinates and timezone:

   - Find coordinates: https://www.latlong.net/
   - Find timezone: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones

3. **Run the installer**

   ```bash
   ./scripts/install.sh
   ```

   This will:

   - Install Python dependencies
   - Calculate today's sunrise/sunset times
   - Set up automatic scheduling with launchd

## Usage

The system runs automatically, but you can also control it manually.

**Real-time updates:** All plugins update running applications immediately (iTerm2) or within 1-2 seconds (Cursor). No restarts needed!

```bash
# Apply mode based on current time (light during day, dark at night)
python3 -m day_night_cycle auto

# Force light mode
python3 -m day_night_cycle light

# Force dark mode
python3 -m day_night_cycle dark

# Show current status and schedule
python3 -m day_night_cycle status

# Show next transition time
python3 -m day_night_cycle next
```

## Configuration

Edit `config.yaml` to customize behavior:

```yaml
location:
  name: "Walla Walla"
  latitude: 46.0645
  longitude: -118.3430
  timezone: "America/Los_Angeles"

plugins:
  # iTerm2 - Use scripts/list_iterm_presets.py to see available presets
  - name: iterm2
    enabled: true
    light_preset: "Light Background" # Configurable
    dark_preset: "Dark Background" # Configurable

  # Claude Code - Switches between 'light' and 'dark' themes
  - name: claude-code
    enabled: true

  # Cursor - Specify exact theme names from Cursor's theme picker
  - name: cursor
    enabled: true
    light_theme: "GitHub Light" # Configurable
    dark_theme: "Dark Modern" # Configurable
```

## Finding Theme Names

Each application uses specific theme/preset names. Use these helper scripts to find the exact names:

**iTerm2 presets:**

```bash
python3 scripts/list_iterm_presets.py
```

**Cursor themes:**

```bash
python3 scripts/list_cursor_themes.py
```

Then update `config.yaml` with the exact names you want to use.

## Creating Custom Plugins

Adding support for new applications is straightforward:

1. **Create a new plugin file** in `day_night_cycle/plugins/`:

   ```python
   from .base import Plugin

   class MyAppPlugin(Plugin):
       @property
       def name(self) -> str:
           return "myapp"

       def set_light_mode(self) -> bool:
           # Implement light mode switching
           return True

       def set_dark_mode(self) -> bool:
           # Implement dark mode switching
           return True
   ```

2. **Import your plugin** in `day_night_cycle/plugins/__init__.py`:

   ```python
   from . import myapp
   ```

3. **Add configuration** in `config.yaml`:

   ```yaml
   plugins:
     - name: myapp
       enabled: true
   ```

## How It Works

1. **Installation**: The installer calculates sunrise/sunset times and creates three launchd agents
2. **Daily updates**: At 12:05 AM, a script recalculates sunrise/sunset times for the new day
3. **Automatic switching**: At sunrise and sunset, the appropriate mode is applied to all enabled plugins
4. **Plugin system**: Each plugin implements a simple interface for setting light and dark modes

## Troubleshooting

**Themes not switching:**

- Check if plugins are enabled in `config.yaml`
- Verify application settings paths exist
- Check logs in `logs/` directory

**Schedule not updating:**

- Verify launchd agents are loaded: `launchctl list | grep daynightcycle`
- Check updater logs: `cat logs/updater.log`

**Manual schedule update:**

```bash
python3 scripts/update_schedule.py
```

## Uninstallation

```bash
./scripts/uninstall.sh
```

This removes launchd agents but keeps the project directory intact.
