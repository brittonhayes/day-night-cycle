# Contributing to day-night-cycle

Thank you for your interest in contributing to day-night-cycle! This document provides guidelines and instructions for contributing.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- A clear, descriptive title
- Steps to reproduce the issue
- Expected behavior vs actual behavior
- Your environment (OS, Go version, architecture, etc.)
- Any relevant logs or error messages

### Suggesting Features

Feature suggestions are welcome! Please create an issue with:
- A clear description of the feature
- The problem it solves or value it adds
- Any relevant examples or use cases

### Adding Support for New Apps

Want to add support for a new application? Great! Here's how:

1. **Create a new plugin function** in `plugins.go`:
   ```go
   func yourAppPlugin(cfg map[string]interface{}, isLight bool) error {
       // Extract config values
       lightTheme, _ := cfg["light_theme"].(string)
       darkTheme, _ := cfg["dark_theme"].(string)

       // Set defaults
       if lightTheme == "" {
           lightTheme = "default-light"
       }
       if darkTheme == "" {
           darkTheme = "default-dark"
       }

       // Determine target theme
       targetTheme := darkTheme
       if isLight {
           targetTheme = lightTheme
       }

       // Implementation here

       return nil
   }
   ```

2. **Register your plugin** in the `plugins` map in `plugins.go`:
   ```go
   var plugins = map[string]PluginFunc{
       // ...
       "your-app": yourAppPlugin,
   }
   ```

3. **Test thoroughly** on your system:
   ```bash
   make build
   ./day-night-cycle light
   ./day-night-cycle dark
   ```

4. **Update documentation** in the README

5. **Submit a pull request** (see below)

### Pull Request Process

1. **Fork the repository** and create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** with clear, focused commits

3. **Test your changes** thoroughly:
   ```bash
   make build
   ./day-night-cycle auto
   ./day-night-cycle light
   ./day-night-cycle dark
   ./day-night-cycle status
   ```

4. **Update documentation** if needed (README, config examples, etc.)

5. **Submit a pull request** with:
   - A clear title describing the change
   - A description of what changed and why
   - Reference to any related issues

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/brittonhayes/day-night-cycle.git
   cd day-night-cycle
   ```

2. Build the project:
   ```bash
   make build
   ```

3. Make your changes and test locally

## Code Style

- Follow standard Go conventions (gofmt)
- Use descriptive variable and function names
- Keep functions focused and single-purpose
- Return descriptive errors using `fmt.Errorf()`
- Add comments for complex logic

## Plugin Guidelines

When creating plugins:
- Return descriptive errors on failure
- Test on the actual application when possible
- Document any special requirements or setup
- Keep configuration simple and intuitive
- Handle missing config files gracefully
- Skip unnecessary writes if already in target mode

## Building and Testing

Run all tests:
```bash
make test
```

Build for multiple platforms:
```bash
make build-all
```

Clean build artifacts:
```bash
make clean
```

## Questions?

Feel free to open an issue with your question or reach out to the maintainers.

## License

By contributing, you agree that your contributions will be released into the public domain under the Unlicense.
