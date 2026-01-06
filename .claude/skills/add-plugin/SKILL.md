---
name: add-plugin
description: Research and create a new plugin for the day/night cycle automation. Use when user asks to "add a plugin for [app name]" or "add support for [app name]".
allowed-tools: WebSearch, WebFetch, Read, Write, Edit, Glob, Grep, Bash
---

# Add Plugin

Automatically research and create a new plugin for the day/night cycle automation system.

## Overview

This skill helps you add support for new applications to the day/night cycle automation. Given an app name, it will:

1. Research how theme switching works for the application
2. Identify configuration file locations or APIs
3. Create a new plugin file in the plugins/ directory
4. Register the plugin in plugins/plugin.go Registry map
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

Create a new plugin file in the `plugins/` directory:

**Read existing plugins first:**
```bash
Read: plugins/iterm2.go  # or any other plugin as example
Read: plugins/plugin.go  # for helper functions
```

**Create new file:** `plugins/[app_name].go`

**Plugin function signature:**
```go
package plugins

func [AppName](cfg map[string]interface{}, isLight bool) error {
    // Implementation
}
```

**Implementation patterns by type:**

1. **JSON Config Files:**
   - Use helper: `UpdateJSONTheme(path, key, value)` from plugins/plugin.go
   - Or read current file into struct or map
   - Update theme property
   - Marshal and write back
   - Verify write succeeded

2. **AppleScript (macOS apps):**
   - Use `exec.Command("osascript", "-e", script)`
   - Set appropriate timeout
   - Handle errors gracefully

3. **CLI Commands:**
   - Use `exec.Command()` with arguments
   - Check error return values
   - Capture and handle stderr

**Available helpers in plugins/plugin.go:**
- `UpdateJSONTheme(path, key, value)` - Update JSON config files
- `ExpandPath(path)` - Expand ~ in file paths

**Error handling requirements:**
- Return descriptive errors using `fmt.Errorf()`
- Handle file not found separately from other errors
- Print helpful messages to stderr when appropriate

### 3. Integration Phase

Register the plugin in the Registry:

**Edit `plugins/plugin.go`:**
```go
var Registry = map[string]Func{
    // ... existing plugins
    "[app-name]": [AppName],
}
```

**Test the plugin:**
```bash
# Build first
make build

# Test light mode
./bin/day-night-cycle --config config.yaml light

# Test dark mode
./bin/day-night-cycle --config config.yaml dark

# Check status
./bin/day-night-cycle --config config.yaml status
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

## Checklist

Before completing the task, verify:

- [ ] New plugin file created in `plugins/` directory (e.g., `plugins/vscode.go`)
- [ ] File has `package plugins` declaration
- [ ] Function signature matches `func [AppName](cfg map[string]interface{}, isLight bool) error`
- [ ] Extracts config values with proper type assertions
- [ ] Handles both light and dark modes based on `isLight` parameter
- [ ] Returns descriptive errors
- [ ] Plugin registered in `plugins/plugin.go` Registry map
- [ ] Configuration example provided
- [ ] Basic testing completed

## Example Usage

User: "Add a plugin for Visual Studio Code"

Expected workflow:
1. Research VS Code settings location and theme configuration
2. Discover settings.json at `~/Library/Application Support/Code/User/settings.json`
3. Find that theme is controlled by `workbench.colorTheme` property
4. Create `plugins/vscode.go` with `VSCode` function
5. Register in `plugins/plugin.go` Registry map as "vscode"
6. Provide config example with popular theme names
7. Test the implementation

## Common Pitfalls

- **Don't guess paths or APIs** - Always research first
- **Check for platform differences** - macOS vs Linux vs Windows
- **Verify theme names** - Use exact names from the application
- **Test error cases** - What if app isn't installed?
- **Consider permissions** - Some apps may need accessibility permissions

## Resources

- Existing plugins: See files in `plugins/` directory for implementation patterns
- Plugin helpers: See `plugins/plugin.go` for available helper functions
- Plugin patterns: See [plugin-patterns.md](plugin-patterns.md) in this skill directory

## Notes

- Always read existing plugin files in plugins/ directory before implementing
- Each plugin gets its own file (e.g., `plugins/iterm2.go`)
- Use helper functions from `plugins/plugin.go` when appropriate
- Follow the patterns from existing plugins
- Research thoroughly before implementing
- Return descriptive errors
- Test both light and dark mode transitions
