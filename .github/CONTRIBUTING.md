# Contributing to day-night-cycle

Thank you for your interest in contributing to day-night-cycle! This document provides guidelines and instructions for contributing.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- A clear, descriptive title
- Steps to reproduce the issue
- Expected behavior vs actual behavior
- Your environment (OS, Python version, etc.)
- Any relevant logs or error messages

### Suggesting Features

Feature suggestions are welcome! Please create an issue with:
- A clear description of the feature
- The problem it solves or value it adds
- Any relevant examples or use cases

### Adding Support for New Apps

Want to add support for a new application? Great! Here's how:

1. **Create a new plugin** in `day_night_cycle/plugins/your_app.py`:
   ```python
   from .base import Plugin

   class YourAppPlugin(Plugin):
       @property
       def name(self) -> str:
           return "your_app"

       def set_light_mode(self) -> bool:
           # Implementation here
           return True

       def set_dark_mode(self) -> bool:
           # Implementation here
           return True
   ```

2. **Import your plugin** in `day_night_cycle/plugins/__init__.py`

3. **Test thoroughly** on your system

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
   python3 -m day_night_cycle auto
   python3 -m day_night_cycle light
   python3 -m day_night_cycle dark
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

2. Install in development mode:
   ```bash
   pip install -e .
   ```

3. Make your changes and test locally

## Code Style

- Follow standard Python conventions (PEP 8)
- Use type hints where appropriate
- Keep functions focused and single-purpose
- Add comments for complex logic

## Plugin Guidelines

When creating plugins:
- Handle errors gracefully (return `False` on failure)
- Test on the actual application when possible
- Document any special requirements or setup
- Keep configuration simple and intuitive

## Questions?

Feel free to open an issue with your question or reach out to the maintainers.

## License

By contributing, you agree that your contributions will be released into the public domain under the Unlicense.
