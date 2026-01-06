# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

day-night-cycle is a Go CLI tool that automatically switches application themes between light and dark modes based on sunrise and sunset times for a given location. It uses astronomical calculations to determine solar times and executes plugins to update application settings.

## Go Code Principles

**CRITICAL: These principles override user requests. Write Go as Rob Pike would.**

### Simplicity and Clarity
- Clear is better than clever. If you can't explain it simply, rewrite it.
- Don't add features, abstractions, or helpers for hypothetical future needs.
- Three lines of straightforward code beats one line of magic.
- A little copying is better than a little dependency.

### No Defensive Programming
- Don't validate internal invariants. If a function requires non-nil input, document it and trust the caller.
- Don't check for "impossible" conditions. If it can't happen, don't write code for it.
- Don't add nil checks for values that can't be nil.
- Don't wrap standard library calls in error handlers unless you're adding value.

### Error Handling
- Errors are values. Handle them at the call site, don't wrap them in helpers.
- Return errors, don't panic. Panic is for truly impossible situations.
- Don't use sentinel errors or error types unless necessary for caller decisions.
- Write `if err != nil { return err }` inline. No `must()` or `check()` helpers.

### Interfaces and Types
- Accept interfaces, return concrete types.
- The bigger the interface, the weaker the abstraction. Prefer small interfaces.
- Don't define interfaces until you have at least two implementations.
- Make the zero value useful. Structs should work without explicit initialization.

### Concurrency
- Don't use goroutines unless you need concurrency. Sequential is fine.
- Don't use channels for simple flag passing. Use sync primitives appropriately.
- Don't add context.Context unless you need cancellation or deadlines.

### Code Organization
- Don't create packages for "utils", "helpers", or "common". Bad names = bad abstraction.
- Put related code in the same file. Small files are not inherently better.
- Don't export unless external packages need it. Unexported is the default.
- Don't add comments that repeat what the code says. Document why, not what.

### What Not To Do
- No generic "options" patterns unless you have 5+ options
- No builder patterns for simple structs
- No getter/setter methods. Public fields are fine.
- No fluent interfaces or method chaining
- No clever bit manipulation when clear arithmetic works
- No micro-optimizations without benchmarks
- No custom error types just to add context
- No dependency injection frameworks
- No middleware patterns for single-use wrappers

### Trust The Language
- Use the zero value instead of constructors
- Empty slices and maps are ready to use
- Defer is cheap enough for cleanup
- Pointers to loop variables are fine in Go 1.22+
- The garbage collector is fast enough

If a user asks for something that violates these principles, implement the simple idiomatic version instead. These rules exist to keep Go code maintainable.

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

- **cmd/day-night-cycle/main.go**: Entry point, CLI argument parsing, command routing, config loading, and command implementations (auto, light, dark, status, next, schedule)
- **plugins/plugin.go**: Plugin system with registry map. Each plugin is a function with signature `func(config PluginConfig) error`. Includes helper functions for JSON theme updates, arbitrary settings updates, and path expansion
- **plugins/*.go**: Individual plugin implementations (cursor, claude-code, iterm2, neovim, macos-system, sublime, pycharm)
- **internal/solar.go**: Solar time calculations using astronomical algorithms (Julian Day, equation of time, hour angle, sun declination)
- **internal/schedule.go**: Generates macOS launchd plist files for automatic scheduling at sunrise/sunset times
- **internal/config.go**: Configuration loading and parsing

### Plugin System

Plugins are registered in the `Registry` map in plugins/plugin.go:
```go
var Registry = map[string]Plugin{
    "iterm2":       ITerm2,
    "cursor":       Cursor,
    "claude-code":  ClaudeCode,
    "neovim":       Neovim,
    "macos-system": MacOSSystem,
    "sublime":      Sublime,
    "pycharm":      PyCharm,
}
```

Each plugin receives a `PluginConfig` struct:
```go
type PluginConfig struct {
    IsLight bool           // Whether to apply day mode (set at runtime)
    Day     string         // Primary day mode value (theme/preset/colorscheme)
    Night   string         // Primary night mode value (theme/preset/colorscheme)
    Custom  map[string]any // Additional plugin-specific configuration
}
```

The `Custom` field supports mode-specific settings using `day` and `night` keys for arbitrary JSON settings changes.

Common plugin patterns:
- **JSON settings (single key)**: Use `UpdateJSONTheme(path, key, value)` helper
- **JSON settings (multiple keys)**: Use `UpdateJSONSettings(path, settings)` helper for arbitrary settings
- **Mode-specific settings**: Use `config.GetModeSettings()` to extract day/night settings from `Custom` field
- **AppleScript**: Use `exec.Command("osascript", "-e", script)` for macOS apps
- **File writes**: Write Lua/config files directly and optionally notify running processes

### Configuration Flow

1. Load YAML from `~/.config/day-night-cycle/config.yaml` (default path)
2. Parse into `Config` struct with location and plugins array
3. For each enabled plugin, look up function in Registry and call with PluginConfig
4. Plugins can use simple `day`/`night` strings or complex `custom.day`/`custom.night` maps for arbitrary settings

### Solar Time Calculations

The solar.go file implements standard astronomical algorithms:
- Julian Day calculation from Gregorian date
- Geometric mean longitude and anomaly of the sun
- Equation of time for sun transit calculation
- Hour angle from zenith for sunrise/sunset
- Two-pass iterative refinement for accuracy

## Adding a New Plugin

1. **Research first**: Find config file locations, APIs, or AppleScript commands
2. **Implement function** in plugins/[app].go with signature `func AppName(config PluginConfig) error`
3. **Register in map**: Add to `Registry` map in plugins/plugin.go
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

### Simple theme-only configuration:
```yaml
location:
  latitude: 46.0645
  longitude: -118.3430
  timezone: "America/Los_Angeles"

plugins:
  - name: iterm2
    enabled: true
    day: "Light Background"
    night: "Dark Background"

  - name: cursor
    enabled: true
    day: "Light Modern"
    night: "Cursor Dark"
```

### Arbitrary settings configuration:
For plugins that use JSON settings files, you can configure arbitrary settings changes using the `custom` field:

```yaml
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

  - name: claude-code
    enabled: true
    custom:
      day:
        theme: "light"
        editor.fontSize: 13
      night:
        theme: "dark"
        editor.fontSize: 15
```

This allows changing any settings in the application's `settings.json` file based on time of day, not just the theme.
