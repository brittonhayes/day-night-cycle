# Example Plugin Implementations

This document shows complete examples of adding plugins for common applications.

## Example 1: Visual Studio Code

**User request:** "Add a plugin for VS Code"

### Research findings:
- Settings location: `~/Library/Application Support/Code/User/settings.json`
- Theme property: `workbench.colorTheme`
- Popular themes: "Light+ (default light)", "Dark+ (default dark)"
- Auto-reload: Yes, VS Code watches settings.json

### Implementation:

Create `plugins/vscode.go`:

```go
package plugins

import (
    "os"
    "path/filepath"
)

// VSCode updates Visual Studio Code theme settings.
func VSCode(cfg map[string]interface{}, isLight bool) error {
    themeKey := "dark_theme"
    defaultTheme := "Dark+ (default dark)"
    if isLight {
        themeKey = "light_theme"
        defaultTheme = "Light+ (default light)"
    }

    theme, ok := cfg[themeKey].(string)
    if !ok {
        theme = defaultTheme
    }

    settingsPath := filepath.Join(
        os.Getenv("HOME"),
        "Library/Application Support/Code/User/settings.json",
    )

    return UpdateJSONTheme(settingsPath, "workbench.colorTheme", theme)
}
```

### Alternate Implementation (Manual):

```go
package plugins

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

// VSCode updates Visual Studio Code theme settings.
func VSCode(cfg map[string]interface{}, isLight bool) error {
    lightTheme, _ := cfg["light_theme"].(string)
    darkTheme, _ := cfg["dark_theme"].(string)

    if lightTheme == "" {
        lightTheme = "Light+ (default light)"
    }
    if darkTheme == "" {
        darkTheme = "Dark+ (default dark)"
    }

    targetTheme := darkTheme
    if isLight {
        targetTheme = lightTheme
    }

    home, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get home directory: %w", err)
    }
    settingsPath := filepath.Join(home, "Library", "Application Support", "Code", "User", "settings.json")

    var settings map[string]interface{}
    data, err := os.ReadFile(settingsPath)
    if err != nil {
        if os.IsNotExist(err) {
            settings = make(map[string]interface{})
        } else {
            return fmt.Errorf("failed to read settings: %w", err)
        }
    } else {
        if err := json.Unmarshal(data, &settings); err != nil {
            return fmt.Errorf("failed to parse settings: %w", err)
        }
    }

    if currentTheme, ok := settings["workbench.colorTheme"].(string); ok && currentTheme == targetTheme {
        return nil
    }

    settings["workbench.colorTheme"] = targetTheme

    updatedData, err := json.MarshalIndent(settings, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal settings: %w", err)
    }

    if err := os.MkdirAll(filepath.Dir(settingsPath), 0755); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }

    if err := os.WriteFile(settingsPath, updatedData, 0644); err != nil {
        return fmt.Errorf("failed to write settings: %w", err)
    }

    return nil
}
```

### Register in plugins/plugin.go:

```go
var Registry = map[string]Func{
    // ... existing plugins
    "vscode": VSCode,
}
```

### Configuration:

Add to `config.yaml`:
```yaml
plugins:
  - name: vscode
    enabled: true
    light_theme: "GitHub Light"
    dark_theme: "GitHub Dark"
```

### Testing:
```bash
make build
./bin/day-night-cycle --config config.yaml light  # Should switch VS Code to light theme
./bin/day-night-cycle --config config.yaml dark   # Should switch VS Code to dark theme
```

## Example 2: Kitty Terminal

**User request:** "Add support for Kitty terminal"

### Research findings:
- Settings location: `~/.config/kitty/kitty.conf`
- Theme approach: Include theme files
- Theme location: `~/.config/kitty/themes/`
- Reload command: `kill -SIGUSR1 $(pgrep kitty)`

### Implementation:

Create `plugins/kitty.go`:

```go
package plugins

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

// Kitty updates Kitty terminal theme configuration.
func Kitty(cfg map[string]interface{}, isLight bool) error {
    lightTheme, _ := cfg["light_theme"].(string)
    darkTheme, _ := cfg["dark_theme"].(string)

    if lightTheme == "" {
        lightTheme = "light"
    }
    if darkTheme == "" {
        darkTheme = "dark"
    }

    themeName := darkTheme
    if isLight {
        themeName = lightTheme
    }

    home, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get home directory: %w", err)
    }
    configPath := filepath.Join(home, ".config", "kitty", "kitty.conf")

    data, err := os.ReadFile(configPath)
    if err != nil {
        return fmt.Errorf("failed to read config: %w", err)
    }

    lines := strings.Split(string(data), "\n")
    newLines := []string{}

    for _, line := range lines {
        if !strings.HasPrefix(strings.TrimSpace(line), "include themes/") {
            newLines = append(newLines, line)
        }
    }

    newLines = append(newLines, fmt.Sprintf("include themes/%s.conf", themeName))

    content := strings.Join(newLines, "\n")
    if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
        return fmt.Errorf("failed to write config: %w", err)
    }

    // Reload Kitty if running
    cmd := exec.Command("pgrep", "kitty")
    if output, err := cmd.Output(); err == nil && len(output) > 0 {
        pid := strings.TrimSpace(string(output))
        exec.Command("kill", "-SIGUSR1", pid).Run()
    }

    return nil
}
```

### Register in plugins/plugin.go:

```go
var Registry = map[string]Func{
    // ... existing plugins
    "kitty": Kitty,
}
```

### Configuration:

```yaml
plugins:
  - name: kitty
    enabled: true
    light_theme: "light"  # Corresponds to themes/light.conf
    dark_theme: "dark"    # Corresponds to themes/dark.conf
```

## Example 3: Alacritty Terminal

**User request:** "Add Alacritty support"

### Research findings:
- Settings location: `~/.config/alacritty/alacritty.yml` or `.toml`
- Theme approach: Import external theme files
- Can use YAML or TOML format

### Implementation:

Create `plugins/alacritty.go`:

```go
package plugins

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "gopkg.in/yaml.v3"
)

// Alacritty updates Alacritty terminal theme configuration.
func Alacritty(cfg map[string]interface{}, isLight bool) error {
    lightTheme, _ := cfg["light_theme"].(string)
    darkTheme, _ := cfg["dark_theme"].(string)

    if lightTheme == "" {
        lightTheme = "themes/light.yml"
    }
    if darkTheme == "" {
        darkTheme = "themes/dark.yml"
    }

    themePath := darkTheme
    if isLight {
        themePath = lightTheme
    }

    home, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get home directory: %w", err)
    }
    configPath := filepath.Join(home, ".config", "alacritty", "alacritty.yml")

    var config map[string]interface{}
    data, err := os.ReadFile(configPath)
    if err != nil {
        return fmt.Errorf("failed to read config: %w", err)
    }

    if err := yaml.Unmarshal(data, &config); err != nil {
        return fmt.Errorf("failed to parse YAML: %w", err)
    }

    // Update import path
    imports, ok := config["import"].([]interface{})
    if !ok {
        imports = []interface{}{}
    }

    // Remove old theme imports
    newImports := []interface{}{}
    for _, imp := range imports {
        if impStr, ok := imp.(string); ok {
            if !strings.Contains(impStr, "themes/") {
                newImports = append(newImports, imp)
            }
        }
    }

    // Add new theme import
    newImports = append(newImports, themePath)
    config["import"] = newImports

    updatedData, err := yaml.Marshal(config)
    if err != nil {
        return fmt.Errorf("failed to marshal YAML: %w", err)
    }

    if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
        return fmt.Errorf("failed to write config: %w", err)
    }

    return nil
}
```

### Configuration:

```yaml
plugins:
  - name: alacritty
    enabled: true
    light_theme: "themes/github-light.yml"
    dark_theme: "themes/github-dark.yml"
```

## Example 4: Sublime Text

**User request:** "Add Sublime Text support"

### Research findings:
- Settings location: `~/Library/Application Support/Sublime Text/Packages/User/Preferences.sublime-settings`
- Theme property: `color_scheme`
- Auto-reload: Yes

### Implementation:

Create `plugins/sublime.go`:

```go
package plugins

import (
    "os"
    "path/filepath"
)

// Sublime updates Sublime Text theme settings.
func Sublime(cfg map[string]interface{}, isLight bool) error {
    themeKey := "dark_theme"
    defaultTheme := "Packages/Color Scheme - Default/Monokai.sublime-color-scheme"
    if isLight {
        themeKey = "light_theme"
        defaultTheme = "Packages/Color Scheme - Default/Breakers.sublime-color-scheme"
    }

    theme, ok := cfg[themeKey].(string)
    if !ok {
        theme = defaultTheme
    }

    settingsPath := filepath.Join(
        os.Getenv("HOME"),
        "Library/Application Support/Sublime Text/Packages/User/Preferences.sublime-settings",
    )

    return UpdateJSONTheme(settingsPath, "color_scheme", theme)
}
```

### Alternate Implementation (Manual):

```go
package plugins

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

// Sublime updates Sublime Text theme settings.
func Sublime(cfg map[string]interface{}, isLight bool) error {
    lightTheme, _ := cfg["light_theme"].(string)
    darkTheme, _ := cfg["dark_theme"].(string)

    if lightTheme == "" {
        lightTheme = "Packages/Color Scheme - Default/Breakers.sublime-color-scheme"
    }
    if darkTheme == "" {
        darkTheme = "Packages/Color Scheme - Default/Monokai.sublime-color-scheme"
    }

    targetTheme := darkTheme
    if isLight {
        targetTheme = lightTheme
    }

    home, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get home directory: %w", err)
    }
    settingsPath := filepath.Join(home, "Library", "Application Support", "Sublime Text", "Packages", "User", "Preferences.sublime-settings")

    var settings map[string]interface{}
    data, err := os.ReadFile(settingsPath)
    if err != nil {
        if os.IsNotExist(err) {
            settings = make(map[string]interface{})
        } else {
            return fmt.Errorf("failed to read settings: %w", err)
        }
    } else {
        if err := json.Unmarshal(data, &settings); err != nil {
            return fmt.Errorf("failed to parse settings: %w", err)
        }
    }

    if currentTheme, ok := settings["color_scheme"].(string); ok && currentTheme == targetTheme {
        return nil
    }

    settings["color_scheme"] = targetTheme

    updatedData, err := json.MarshalIndent(settings, "", "    ")
    if err != nil {
        return fmt.Errorf("failed to marshal settings: %w", err)
    }

    if err := os.WriteFile(settingsPath, updatedData, 0644); err != nil {
        return fmt.Errorf("failed to write settings: %w", err)
    }

    return nil
}
```

### Configuration:

```yaml
plugins:
  - name: sublime
    enabled: true
    light_theme: "Packages/Theme - GitHub/GitHub-light.tmTheme"
    dark_theme: "Packages/Theme - GitHub/GitHub-dark.tmTheme"
```

## Common Patterns Observed

1. **JSON config files** (VS Code, Cursor, Electron apps)
   - Most reliable method
   - Usually auto-reload
   - Need exact theme names
   - **Use `UpdateJSONTheme()` helper from plugins/plugin.go**

2. **Config file includes** (Kitty, Alacritty)
   - Include external theme files
   - May need manual reload signal

3. **Property-based configs** (.vimrc, .bashrc)
   - Update specific lines
   - Usually need app restart

4. **State files** (some apps)
   - JSON but not traditional settings
   - Often require app restart

## Tips for Quick Implementation

1. **Create separate file** - Each plugin gets its own file in `plugins/`
2. **Start with research** - Don't guess file locations
3. **Check existing plugins** - Follow established patterns in `plugins/` directory
4. **Use helpers** - `UpdateJSONTheme()` and `ExpandPath()` from `plugins/plugin.go`
5. **Test error cases** - What if file doesn't exist?
6. **Verify writes** - Check if the change actually saved
7. **Note reload requirements** - Does app need restart?
8. **Use exact theme names** - Copy from app's UI
9. **Handle platform differences** - Windows/Linux/macOS paths differ
10. **Return descriptive errors** - Guide users to fix issues
11. **Register in Registry** - Add to `plugins/plugin.go` Registry map
