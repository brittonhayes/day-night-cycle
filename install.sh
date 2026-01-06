#!/bin/bash
# Standalone installation script for day-night-cycle
# Can be run directly via: curl -fsSL https://raw.githubusercontent.com/brittonhayes/day-night-cycle/main/install.sh | bash

set -e

INSTALL_DIR="$HOME/.config/day-night-cycle"
REPO_URL="https://github.com/brittonhayes/day-night-cycle.git"
LAUNCHD_DIR="$HOME/Library/LaunchAgents"

echo "=========================================="
echo "day-night-cycle installation"
echo "=========================================="
echo ""

# Check for required commands
if ! command -v python3 &>/dev/null; then
  echo "Error: Python 3 is required but not found"
  echo "Install Python 3 from https://www.python.org/downloads/"
  exit 1
fi

if ! command -v git &>/dev/null; then
  echo "Error: git is required but not found"
  echo "Install git from https://git-scm.com/downloads"
  exit 1
fi

PYTHON_PATH=$(which python3)
echo "Found Python: $PYTHON_PATH"
echo ""

# Clone or update repository
if [ -d "$INSTALL_DIR" ]; then
  echo "Updating existing installation at $INSTALL_DIR..."
  cd "$INSTALL_DIR"
  git pull --quiet
else
  echo "Cloning repository to $INSTALL_DIR..."
  git clone --quiet "$REPO_URL" "$INSTALL_DIR"
  cd "$INSTALL_DIR"
fi

echo "Installing Python dependencies..."
pip3 install -r requirements.txt --quiet --user

# Create logs directory
mkdir -p "$INSTALL_DIR/logs"

# Interactive configuration
echo ""
echo "=========================================="
echo "Configuration"
echo "=========================================="
echo ""

if [ -f "$INSTALL_DIR/config.yaml" ]; then
  echo "Found existing config.yaml"
  read -p "Do you want to reconfigure? (y/N): " reconfigure
  if [[ ! $reconfigure =~ ^[Yy]$ ]]; then
    echo "Using existing configuration"
  else
    rm "$INSTALL_DIR/config.yaml"
  fi
fi

if [ ! -f "$INSTALL_DIR/config.yaml" ]; then
  echo "Setting up your location..."
  echo ""
  echo "Find your coordinates: https://www.latlong.net/"
  echo "Find your timezone: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones"
  echo ""

  read -p "Enter your latitude (e.g., 46.0645): " latitude
  read -p "Enter your longitude (e.g., -118.3430): " longitude
  read -p "Enter your timezone (e.g., America/Los_Angeles): " timezone

  # Create config from template
  cat >"$INSTALL_DIR/config.yaml" <<EOF
location:
  name: "User Location"
  latitude: $latitude
  longitude: $longitude
  timezone: "$timezone"

plugins:
  - name: iterm2
    enabled: true
    light_preset: "Light Background"
    dark_preset: "Dark Background"

  - name: claude-code
    enabled: true

  - name: cursor
    enabled: true
    light_theme: "GitHub Light"
    dark_theme: "Dark Modern"

  - name: neovim
    enabled: true

  - name: macos-system
    enabled: false
EOF

  echo ""
  echo "Configuration saved to $INSTALL_DIR/config.yaml"
  echo "You can edit this file later to customize plugin settings"
fi

# Create launchd directory if needed
mkdir -p "$LAUNCHD_DIR"

# Create updater plist (runs daily at 12:05 AM)
UPDATER_PLIST="$LAUNCHD_DIR/com.daynightcycle.updater.plist"
cat >"$UPDATER_PLIST" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.daynightcycle.updater</string>
    <key>ProgramArguments</key>
    <array>
        <string>$PYTHON_PATH</string>
        <string>$INSTALL_DIR/scripts/update_schedule.py</string>
    </array>
    <key>WorkingDirectory</key>
    <string>$INSTALL_DIR</string>
    <key>StartCalendarInterval</key>
    <dict>
        <key>Hour</key>
        <integer>0</integer>
        <key>Minute</key>
        <integer>5</integer>
    </dict>
    <key>StandardOutPath</key>
    <string>$INSTALL_DIR/logs/updater.log</string>
    <key>StandardErrorPath</key>
    <string>$INSTALL_DIR/logs/updater.error.log</string>
</dict>
</plist>
EOF

# Calculate today's sunrise/sunset and create schedule
echo ""
echo "Calculating today's sunrise and sunset times..."
python3 "$INSTALL_DIR/scripts/update_schedule.py"

# Load launchd agents
echo ""
echo "Configuring automatic schedule..."

# Unload existing agents (ignore errors)
launchctl unload "$UPDATER_PLIST" 2>/dev/null || true
launchctl unload "$LAUNCHD_DIR/com.daynightcycle.sunrise.plist" 2>/dev/null || true
launchctl unload "$LAUNCHD_DIR/com.daynightcycle.sunset.plist" 2>/dev/null || true

# Load agents
launchctl load "$UPDATER_PLIST"
launchctl load "$LAUNCHD_DIR/com.daynightcycle.sunrise.plist"
launchctl load "$LAUNCHD_DIR/com.daynightcycle.sunset.plist"

# Add convenience alias suggestion
SHELL_CONFIG=""
if [ -n "$ZSH_VERSION" ]; then
  SHELL_CONFIG="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ]; then
  SHELL_CONFIG="$HOME/.bashrc"
fi

echo ""
echo "=========================================="
echo "Installation complete!"
echo "=========================================="
echo ""
echo "Themes will automatically switch at sunrise and sunset."
echo ""
echo "Manual commands:"
echo "  python3 -m day_night_cycle auto     # Apply mode based on time"
echo "  python3 -m day_night_cycle light    # Force light mode"
echo "  python3 -m day_night_cycle dark     # Force dark mode"
echo "  python3 -m day_night_cycle status   # Show status"
echo ""
echo "Configuration: $INSTALL_DIR/config.yaml"
echo ""

if [ -n "$SHELL_CONFIG" ]; then
  echo "Optional: Add an alias to your $SHELL_CONFIG:"
  echo "  alias dnc='cd $INSTALL_DIR && python3 -m day_night_cycle'"
  echo ""
fi

echo "To uninstall:"
echo "  bash $INSTALL_DIR/scripts/uninstall.sh"
echo ""
