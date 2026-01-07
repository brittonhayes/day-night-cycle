# day-night-cycle

Switch application themes at sunrise and sunset.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/brittonhayes/day-night-cycle/main/install.sh | bash
```

Or:

```bash
go install github.com/brittonhayes/day-night-cycle@latest
```

Or download a binary from the [releases page](https://github.com/brittonhayes/day-night-cycle/releases).

## Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/brittonhayes/day-night-cycle/main/install.sh | bash -s -- --uninstall
```

This will:

- Unload the launchd agent
- Remove `/usr/local/bin/day-night-cycle`
- Remove `~/.config/day-night-cycle/` (configuration directory)
- Remove `~/Library/LaunchAgents/com.daynightcycle.schedule.plist`

## Supported Plugins

- **iterm2** - iTerm2 terminal
- **cursor** - Cursor editor (supports arbitrary settings)
- **claude-code** - Claude Code editor (supports arbitrary settings)
- **neovim** - Neovim editor
- **macos-system** - macOS system appearance
- **sublime** - Sublime Text editor (supports arbitrary settings)
- **pycharm** - PyCharm IDE (supports arbitrary settings)

## Configure

Edit `~/.config/day-night-cycle/config.yaml`:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/brittonhayes/day-night-cycle/main/config.schema.json
location:
  latitude: 46.0645
  longitude: -118.3430
  timezone: "America/Los_Angeles"
  dayOffset: "1h"

plugins:
  - name: iterm2
    enabled: true
    day: "Light Background"
    night: "Dark Background"

  - name: cursor
    enabled: true
    day: "Light Modern"
    night: "Cursor Dark"

  - name: claude-code
    enabled: true

  - name: neovim
    enabled: true
    day: "github_light"
    night: "github_dark_default"

  - name: macos-system
    enabled: false
```

### Arbitrary Settings

For plugins that use JSON settings files (like Cursor and Claude Code), you can configure arbitrary settings changes using the `custom` field:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/brittonhayes/day-night-cycle/main/config.schema.json
plugins:
  - name: cursor
    enabled: true
    custom:
      day:
        workbench.colorTheme: "Default Light+"
        editor.fontSize: 14
        zenMode.fullScreen: false
      night:
        workbench.colorTheme: "Default Dark+"
        editor.fontSize: 16
        zenMode.fullScreen: true
        editor.lineHeight: 1.6
```

This allows you to change any settings in the application's `settings.json` file based on the time of day, not just the theme.

## Use

```bash
day-night-cycle auto      # apply mode for current time
day-night-cycle light     # force light mode
day-night-cycle dark      # force dark mode
day-night-cycle status    # show current status
day-night-cycle next      # show next transition
day-night-cycle schedule  # generate launchd schedule
```

## Build

```bash
git tag vx.x.x
make build-all
make release
```
