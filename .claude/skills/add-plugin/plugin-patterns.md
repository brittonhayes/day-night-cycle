# Plugin Implementation Patterns

This guide provides detailed patterns for implementing plugins based on how the target application handles theme configuration.

## Pattern 1: JSON Configuration Files

**Use when:** Application stores settings in JSON format (e.g., VS Code, Cursor, Claude Code)

**Example applications:** VS Code, Cursor, Sublime Text, many Electron apps

### Template

```go
func [appName]Plugin(cfg map[string]interface{}, isLight bool) error {
    // Extract config values
    lightTheme, _ := cfg["light_theme"].(string)
    darkTheme, _ := cfg["dark_theme"].(string)

    // Set defaults if not provided
    if lightTheme == "" {
        lightTheme = "[Default Light Theme]"
    }
    if darkTheme == "" {
        darkTheme = "[Default Dark Theme]"
    }

    // Determine target theme
    targetTheme := darkTheme
    if isLight {
        targetTheme = lightTheme
    }

    // Build settings path
    home, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get home directory: %w", err)
    }
    settingsPath := filepath.Join(home, "Library", "Application Support", "[AppName]", "User", "settings.json")

    // Read current settings
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

    // Check if already set
    if currentTheme, ok := settings["[theme_property_name]"].(string); ok && currentTheme == targetTheme {
        return nil // Already set, nothing to do
    }

    // Update theme
    settings["[theme_property_name]"] = targetTheme

    // Marshal and write
    updatedData, err := json.MarshalIndent(settings, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal settings: %w", err)
    }

    // Ensure directory exists
    if err := os.MkdirAll(filepath.Dir(settingsPath), 0755); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }

    // Write file
    if err := os.WriteFile(settingsPath, updatedData, 0644); err != nil {
        return fmt.Errorf("failed to write settings: %w", err)
    }

    return nil
}
```

### Key considerations:
- Use `os.IsNotExist()` to check for missing files
- Handle JSON unmarshal errors gracefully
- Skip writes if already in target mode
- Use `MarshalIndent` for readable output
- Ensure parent directories exist before writing

## Pattern 2: YAML Configuration Files

**Use when:** Application uses YAML for settings

**Example applications:** Many CLI tools, development tools

### Template

```go
func [appName]Plugin(cfg map[string]interface{}, isLight bool) error {
    lightTheme, _ := cfg["light_theme"].(string)
    darkTheme, _ := cfg["dark_theme"].(string)

    if lightTheme == "" {
        lightTheme = "[default-light]"
    }
    if darkTheme == "" {
        darkTheme = "[default-dark]"
    }

    targetTheme := darkTheme
    if isLight {
        targetTheme = lightTheme
    }

    home, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get home directory: %w", err)
    }
    configPath := filepath.Join(home, ".[app-name]", "config.yml")

    // Read current config
    var config map[string]interface{}
    data, err := os.ReadFile(configPath)
    if err != nil {
        if os.IsNotExist(err) {
            config = make(map[string]interface{})
        } else {
            return fmt.Errorf("failed to read config: %w", err)
        }
    } else {
        if err := yaml.Unmarshal(data, &config); err != nil {
            return fmt.Errorf("failed to parse YAML: %w", err)
        }
    }

    // Check if already set
    if currentTheme, ok := config["theme"].(string); ok && currentTheme == targetTheme {
        return nil
    }

    // Update theme
    config["theme"] = targetTheme

    // Marshal and write
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

## Pattern 3: AppleScript Control (macOS)

**Use when:** macOS application supports AppleScript automation

**Example applications:** iTerm2, Terminal, some native macOS apps

### Template

```go
func [appName]Plugin(cfg map[string]interface{}, isLight bool) error {
    lightPreset, _ := cfg["light_preset"].(string)
    darkPreset, _ := cfg["dark_preset"].(string)

    if lightPreset == "" {
        lightPreset = "[Default Light]"
    }
    if darkPreset == "" {
        darkPreset = "[Default Dark]"
    }

    preset := darkPreset
    if isLight {
        preset = lightPreset
    }

    script := fmt.Sprintf(`
        tell application "[App Name]"
            -- AppleScript commands to set theme
            -- Example: set current theme to "%s"
        end tell
    `, preset)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, "osascript", "-e", script)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("AppleScript failed: %w: %s", err, string(output))
    }

    return nil
}
```

### Key considerations:
- Use context with timeout to prevent hanging
- Capture stderr for error messages
- Handle case where osascript isn't available
- Some apps require them to be running

## Pattern 4: Command-Line Interface

**Use when:** Application provides CLI for theme control

**Example applications:** Many terminal apps, CLI-first tools

### Template

```go
func [appName]Plugin(cfg map[string]interface{}, isLight bool) error {
    lightTheme, _ := cfg["light_theme"].(string)
    darkTheme, _ := cfg["dark_theme"].(string)

    if lightTheme == "" {
        lightTheme = "[default-light]"
    }
    if darkTheme == "" {
        darkTheme = "[default-dark]"
    }

    theme := darkTheme
    if isLight {
        theme = lightTheme
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, "[command-name]", "theme", "set", theme)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("command failed: %w: %s", err, string(output))
    }

    return nil
}
```

## Pattern 5: Text File Manipulation

**Use when:** Need to modify specific lines in config files (like .vimrc)

### Template

```go
func [appName]Plugin(cfg map[string]interface{}, isLight bool) error {
    lightSetting, _ := cfg["light_setting"].(string)
    darkSetting, _ := cfg["dark_setting"].(string)

    setting := darkSetting
    if isLight {
        setting = lightSetting
    }

    home, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get home directory: %w", err)
    }
    configPath := filepath.Join(home, ".[app-config]")

    // Read file
    data, err := os.ReadFile(configPath)
    if err != nil && !os.IsNotExist(err) {
        return fmt.Errorf("failed to read config: %w", err)
    }

    lines := []string{}
    if len(data) > 0 {
        lines = strings.Split(string(data), "\n")
    }

    // Remove old setting lines
    newLines := []string{}
    for _, line := range lines {
        if !strings.HasPrefix(strings.TrimSpace(line), "[setting-prefix]") {
            newLines = append(newLines, line)
        }
    }

    // Add new setting
    newLines = append(newLines, fmt.Sprintf("[setting-prefix] %s", setting))

    // Write back
    content := strings.Join(newLines, "\n")
    if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
        return fmt.Errorf("failed to write config: %w", err)
    }

    return nil
}
```

## Common Helper Patterns

### Check if Application is Running

```go
func isAppRunning(appName string) bool {
    cmd := exec.Command("pgrep", "-x", appName)
    return cmd.Run() == nil
}
```

### Send Reload Signal (macOS/Linux)

```go
func reloadApp(appName string) error {
    cmd := exec.Command("pgrep", appName)
    output, err := cmd.Output()
    if err != nil {
        return nil // Not running, nothing to reload
    }

    pid := strings.TrimSpace(string(output))
    if pid == "" {
        return nil
    }

    cmd = exec.Command("kill", "-SIGUSR1", pid)
    return cmd.Run()
}
```

## Research Checklist

When researching a new application, find:

- [ ] Settings file location and format
- [ ] Theme property name(s)
- [ ] Available theme names/identifiers
- [ ] Does app auto-reload settings?
- [ ] Are there platform differences? (macOS/Linux/Windows)
- [ ] Any special permissions needed?
- [ ] Official documentation on automation
- [ ] Existing automation examples in community

## Testing Checklist

After implementation:

- [ ] Plugin function compiles without errors
- [ ] Registered in plugins map
- [ ] Light mode switches successfully
- [ ] Dark mode switches successfully
- [ ] Error messages are helpful
- [ ] Returns proper error values
- [ ] Handles missing files gracefully
- [ ] Works when app is not running (if applicable)
- [ ] Works when app is running (if applicable)
