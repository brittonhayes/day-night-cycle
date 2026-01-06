# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

day-night-cycle is a Go CLI tool that automatically switches application themes between light and dark modes based on sunrise and sunset times for a given location. It uses astronomical calculations to determine solar times and executes plugins to update application settings.

## Commands

### Building and Testing
```bash
# Build for current platform
make build

# Build for all platforms (Intel and Apple Silicon)
make build-all

# Run tests
make test

# Format code
make fmt

# Run go vet
make vet

# Clean build artifacts
make clean
```

### Running the Application
```bash
# Apply mode based on current time
./bin/day-night-cycle auto

# Force light or dark mode
./bin/day-night-cycle light
./bin/day-night-cycle dark

# Show status and schedule
./bin/day-night-cycle status

# Show next transition time
./bin/day-night-cycle next

# Generate launchd schedule
./bin/day-night-cycle schedule
```

### Installation Testing
```bash
# Install to /usr/local/bin
make install

# Test with config file
./bin/day-night-cycle --config ~/.config/day-night-cycle/config.yaml auto
```

## Architecture

### Core Files Structure

- **main.go**: Entry point, CLI argument parsing, command routing, config loading, and command implementations (auto, light, dark, status, next, schedule)
- **plugins.go**: Plugin system with registry map. Each plugin is a function with signature `func(cfg map[string]interface{}, isLight bool) error`. Includes helper functions for JSON theme updates and path expansion
- **solar.go**: Solar time calculations using astronomical algorithms (Julian Day, equation of time, hour angle, sun declination)
- **schedule.go**: Generates macOS launchd plist files for automatic scheduling at sunrise/sunset times

### Plugin System

Plugins are registered in a global `plugins` map in plugins.go:
```go
var plugins = map[string]PluginFunc{
    "iterm2":       iterm2Plugin,
    "cursor":       cursorPlugin,
    "claude-code":  claudeCodePlugin,
    "neovim":       neovimPlugin,
    "macos-system": macosSystemPlugin,
}
```

Each plugin receives:
- `cfg map[string]interface{}`: Plugin-specific config values from YAML
- `isLight bool`: true for light mode, false for dark mode

Common plugin patterns:
- **JSON settings**: Use `updateJSONTheme()` helper for JSON config files
- **AppleScript**: Use `exec.Command("osascript", "-e", script)` for macOS apps
- **File writes**: Write Lua/config files directly and optionally notify running processes

### Configuration Flow

1. Load YAML from `~/.config/day-night-cycle/config.yaml` (default path)
2. Parse into `Config` struct with location and plugins array
3. For each enabled plugin, look up function in registry and call with plugin config
4. Each plugin extracts its specific config values using type assertions

### Solar Time Calculations

The solar.go file implements standard astronomical algorithms:
- Julian Day calculation from Gregorian date
- Geometric mean longitude and anomaly of the sun
- Equation of time for sun transit calculation
- Hour angle from zenith for sunrise/sunset
- Two-pass iterative refinement for accuracy

## Adding a New Plugin

1. **Research first**: Find config file locations, APIs, or AppleScript commands
2. **Implement function** in plugins.go with signature `func [appName]Plugin(cfg map[string]interface{}, isLight bool) error`
3. **Register in map**: Add to `plugins` map in plugins.go
4. **Test thoroughly**: Build and test both light and dark modes
5. **Use the /add-plugin skill** for guided plugin creation

See CONTRIBUTING.md and `.claude/skills/add-plugin/SKILL.md` for detailed plugin development guidance.

## Release Process

```bash
# Tag the release
git tag v1.0.0
git push origin v1.0.0

# Build and create GitHub release (requires gh CLI)
make release
```

The release target builds for both architectures and creates a GitHub release with binaries.

## Configuration Example

Config is YAML at `~/.config/day-night-cycle/config.yaml`:
```yaml
location:
  latitude: 46.0645
  longitude: -118.3430
  timezone: "America/Los_Angeles"

plugins:
  - name: iterm2
    enabled: true
    light_preset: "Light Background"
    dark_preset: "Dark Background"
```

Each plugin's config is passed as a map to its function.
