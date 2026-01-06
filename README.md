# day-night-cycle

Automatically switch your apps between light and dark themes based on sunrise and sunset for your location.

Built with Go for simplicity and performance - single binary, no dependencies.

## Install

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/brittonhayes/day-night-cycle/main/install-go.sh | bash
```

### Alternative: Using Go

```bash
go install github.com/brittonhayes/day-night-cycle@latest
```

### Alternative: Manual Download

Download the appropriate binary for your system from the [releases page](https://github.com/brittonhayes/day-night-cycle/releases):
- `day-night-cycle-darwin-arm64` for Apple Silicon Macs
- `day-night-cycle-darwin-amd64` for Intel Macs

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
day-night-cycle auto

# Force light or dark mode
day-night-cycle light
day-night-cycle dark

# Show status
day-night-cycle status

# Show next transition
day-night-cycle next

# Generate/update launchd schedule
day-night-cycle schedule
```

If you used the install script, you can also use the full path:
```bash
~/.config/day-night-cycle/day-night-cycle auto
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

  - name: claude-code
    enabled: true

  - name: neovim
    enabled: true
    light_colorscheme: "github_light"
    dark_colorscheme: "github_dark_default"

  - name: macos-system
    enabled: false
    # Optional wallpaper switching:
    # light_wallpaper: "~/Pictures/light.jpg"
    # dark_wallpaper: "~/Pictures/dark.jpg"
```

Find coordinates at [latlong.net](https://www.latlong.net/) and timezones at [Wikipedia](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).

## Building from Source

Requirements: Go 1.21 or later

```bash
# Clone the repository
git clone https://github.com/brittonhayes/day-night-cycle.git
cd day-night-cycle

# Build for current platform
make build

# Build for all platforms
make build-all

# Install to /usr/local/bin
make install
```

### Makefile Targets

- `make build` - Build for current platform
- `make build-darwin-amd64` - Build for Intel Mac
- `make build-darwin-arm64` - Build for Apple Silicon
- `make build-all` - Build for both architectures
- `make install` - Install to /usr/local/bin
- `make release` - Create GitHub release (requires `gh` CLI)
- `make clean` - Remove built binaries
- `make test` - Run tests
- `make help` - Show all targets

## Creating a Release

To cut a new release:

```bash
# Tag the release
git tag v1.0.0
git push origin v1.0.0

# Build and create GitHub release
make release
```

This requires the GitHub CLI (`gh`) to be installed and authenticated.

## Uninstall

```bash
# Stop and remove launchd agent
launchctl unload ~/Library/LaunchAgents/com.daynightcycle.schedule.plist
rm ~/Library/LaunchAgents/com.daynightcycle.schedule.plist

# Remove configuration and binary
rm -rf ~/.config/day-night-cycle
```


## License

This project is released into the public domain under the [Unlicense](LICENSE).
