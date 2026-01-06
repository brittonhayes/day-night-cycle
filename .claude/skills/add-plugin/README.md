# Add Plugin Skill

A Claude Code skill for automatically researching and creating new plugins for the day/night cycle automation system.

## Purpose

This skill automates the process of adding new application support to the day/night cycle automation. When a user asks to add support for a new application, this skill will:

1. Research how theme switching works for the target application
2. Identify the appropriate implementation pattern
3. Create a plugin file implementing the Plugin base class
4. Update the plugin registry
5. Provide configuration examples and documentation

## Usage

Simply ask Claude Code to add support for an application:

```
Add a plugin for Visual Studio Code
```

```
Add support for Slack
```

```
Create a plugin for Kitty terminal
```

The skill will automatically activate based on these types of requests.

## What It Does

### Research Phase
- Searches for official documentation on theme switching
- Identifies configuration file locations and formats
- Discovers available theme names
- Determines if live reloading is supported
- Checks for AppleScript or CLI support

### Implementation Phase
- Selects the appropriate plugin pattern (JSON, YAML, AppleScript, CLI, plist)
- Creates a new plugin file in `day_night_cycle/plugins/`
- Implements the required methods: `name`, `set_light_mode()`, `set_dark_mode()`
- Adds proper error handling and validation
- Updates `day_night_cycle/plugins/__init__.py` to register the plugin

### Documentation Phase
- Provides configuration examples for `config.yaml`
- Documents any special setup requirements
- Notes whether app restart or reload is needed
- Suggests theme names to use

## Files

- **SKILL.md** - Main skill definition and instructions
- **plugin-patterns.md** - Detailed implementation patterns by configuration type
- **examples.md** - Complete worked examples for common applications

## Supported Patterns

The skill knows how to implement plugins for:

1. **JSON configuration files** (VS Code, Cursor, Electron apps)
2. **YAML configuration files** (CLI tools, development tools)
3. **AppleScript control** (iTerm2, Terminal, native macOS apps)
4. **Command-line interfaces** (CLI-first tools)
5. **Property lists (plist)** (Native macOS applications)
6. **Hybrid approaches** (apps with multiple control methods)

## Prerequisites

The skill expects:
- The base Plugin class at `day_night_cycle/plugins/base.py`
- Plugin registry at `day_night_cycle/plugins/__init__.py`
- Configuration file at `config.yaml` (or `config.example.yaml`)
- Access to web search for research

## Examples

### Example 1: Adding VS Code Support

**User:** "Add a plugin for Visual Studio Code"

**Skill actions:**
1. Searches for VS Code theme documentation
2. Discovers settings.json location and workbench.colorTheme property
3. Creates `day_night_cycle/plugins/vscode.py` using JSON pattern
4. Updates plugin registry
5. Provides config with popular theme names

### Example 2: Adding Kitty Terminal

**User:** "Add support for Kitty terminal"

**Skill actions:**
1. Researches Kitty configuration
2. Finds it uses include-based theme system
3. Creates plugin that updates kitty.conf
4. Implements reload signal (SIGUSR1) for live updates
5. Documents theme file requirements

## Benefits

- **Saves time**: Automates repetitive research and boilerplate
- **Consistency**: Follows established patterns and conventions
- **Quality**: Includes proper error handling and validation
- **Documentation**: Provides complete setup instructions
- **Best practices**: Uses proven patterns from existing plugins

## Customization

The skill can be extended with:
- Additional plugin patterns
- Platform-specific implementations (Windows, Linux)
- Integration with helper script generation
- Automated testing setup

## Maintenance

To update the skill:
1. Edit the relevant `.md` files in `.claude/skills/add-plugin/`
2. Keep patterns in sync with base Plugin class changes
3. Add new examples as more plugins are implemented
4. Update research checklist based on common issues

## Notes

- The skill uses web search to ensure accurate, up-to-date information
- It follows the patterns established by existing plugins
- All generated code includes proper error handling
- Plugins are tested for basic functionality before completion

## Related Documentation

- [PLUGINS.md](../../../PLUGINS.md) - Plugin development guide
- [day_night_cycle/plugins/base.py](../../../day_night_cycle/plugins/base.py) - Base Plugin class
- [README.md](../../../README.md) - Project overview
