---
name: add-plugin
description: Research and create a new plugin for the day/night cycle automation. Use when user asks to "add a plugin for [app name]" or "add support for [app name]".
allowed-tools: WebSearch, WebFetch, Read, Write, Edit, Glob, Grep, Bash(python:*)
---

# Add Plugin

Automatically research and create a new plugin for the day/night cycle automation system.

## Overview

This skill helps you add support for new applications to the day/night cycle automation. Given an app name, it will:

1. Research how theme switching works for the application
2. Identify configuration file locations or APIs
3. Create a plugin implementing the Plugin base class
4. Update the plugin registry
5. Provide configuration examples

## Instructions

When adding a plugin, follow these steps:

### 1. Research Phase

First, understand how theme switching works for the target application:

- Search for official documentation on theme/appearance switching
- Look for configuration file locations (JSON, YAML, plist, etc.)
- Check for AppleScript support on macOS
- Investigate CLI commands or APIs
- Review GitHub issues or community discussions about automation

**Key questions to answer:**
- Where are theme settings stored?
- What file format is used?
- What are the exact theme/preset names?
- Does the app support live reloading of settings?
- Are there any special requirements or permissions needed?

### 2. Implementation Phase

Create the plugin following the established patterns:

**Read the base Plugin class first:**
```python
# Always start by reading the base class
Read: day_night_cycle/plugins/base.py
```

**Study existing plugins for patterns:**
- `day_night_cycle/plugins/claude_code.py` - JSON file-based config
- `day_night_cycle/plugins/cursor.py` - JSON with verification
- `day_night_cycle/plugins/iterm2.py` - AppleScript-based control

**Create the plugin file:**
- File name: `day_night_cycle/plugins/[app_name].py` (use snake_case)
- Class name: `[AppName]Plugin` (use PascalCase)
- Implement required methods: `name`, `set_light_mode()`, `set_dark_mode()`
- Optional: `validate_config()` for startup validation

**Implementation patterns by type:**

1. **JSON/YAML Config Files:**
   - Read current file
   - Update theme property
   - Write back with proper formatting
   - Verify write succeeded

2. **AppleScript (macOS apps):**
   - Use `subprocess.run(['osascript', '-e', script])`
   - Set appropriate timeout
   - Handle errors gracefully

3. **CLI Commands:**
   - Use `subprocess.run()` with command array
   - Check return codes
   - Capture and handle errors

**Error handling requirements:**
- Always use try/except blocks
- Return `True` on success, `False` on failure
- Print helpful error messages with "    Error: " prefix
- Handle FileNotFoundError separately from generic errors

### 3. Integration Phase

Update the plugin registry:

**Edit `day_night_cycle/plugins/__init__.py`:**
```python
# Add import for new plugin
from . import [app_name]
```

**Test the plugin:**
```bash
# Validate configuration
python3 scripts/test_config.py

# Test light mode
python3 -m day_night_cycle light

# Test dark mode
python3 -m day_night_cycle dark

# Check status
python3 -m day_night_cycle status
```

### 4. Documentation Phase

Provide the user with:

1. **Configuration example** for `config.yaml`:
```yaml
plugins:
  - name: [app-name]
    enabled: true
    # Include any custom options discovered
```

2. **Setup instructions** if needed:
   - Where to find theme names
   - Any prerequisites or permissions
   - How to verify it's working

3. **Helper script** (optional) for listing available themes:
```python
# Example: scripts/list_[app_name]_themes.py
```

## Checklist

Before completing the task, verify:

- [ ] Plugin file created in `day_night_cycle/plugins/`
- [ ] Plugin inherits from `Plugin` base class
- [ ] `name` property returns correct identifier
- [ ] `set_light_mode()` implemented and returns bool
- [ ] `set_dark_mode()` implemented and returns bool
- [ ] Error handling with try/except blocks
- [ ] Plugin imported in `plugins/__init__.py`
- [ ] Configuration example provided
- [ ] Basic testing completed

## Example Usage

User: "Add a plugin for Visual Studio Code"

Expected workflow:
1. Research VS Code settings location and theme configuration
2. Discover settings.json at `~/Library/Application Support/Code/User/settings.json`
3. Find that theme is controlled by `workbench.colorTheme` property
4. Create `day_night_cycle/plugins/vscode.py` following JSON config pattern
5. Update `plugins/__init__.py` to import vscode
6. Provide config example with popular theme names
7. Test the implementation

## Common Pitfalls

- **Don't guess paths or APIs** - Always research first
- **Check for platform differences** - macOS vs Linux vs Windows
- **Verify theme names** - Use exact names from the application
- **Test error cases** - What if app isn't installed?
- **Consider permissions** - Some apps may need accessibility permissions

## Resources

- Plugin development guide: [PLUGINS.md](../../../PLUGINS.md)
- Base Plugin class: [day_night_cycle/plugins/base.py](../../../day_night_cycle/plugins/base.py)
- Existing plugins: [day_night_cycle/plugins/](../../../day_night_cycle/plugins/)
- Plugin patterns: See [plugin-patterns.md](plugin-patterns.md) in this skill directory

## Notes

- Always read the base Plugin class before implementing
- Follow the patterns from existing plugins
- Research thoroughly before implementing
- Provide clear error messages
- Test both light and dark mode transitions
