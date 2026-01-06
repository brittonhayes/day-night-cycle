#!/bin/bash
# Installation script for Day/Night Cycle automation

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LAUNCHD_DIR="$HOME/Library/LaunchAgents"
LOGS_DIR="$PROJECT_DIR/logs"

echo "================================================"
echo "Day/Night Cycle Automation - Installation"
echo "================================================"
echo ""

# Check if Python 3 is available
if ! command -v python3 &> /dev/null; then
    echo "Error: Python 3 is required but not found"
    exit 1
fi

PYTHON_PATH=$(which python3)
echo "Using Python: $PYTHON_PATH"
echo ""

# Create logs directory
echo "Creating logs directory..."
mkdir -p "$LOGS_DIR"

# Install Python dependencies
echo "Installing Python dependencies..."
cd "$PROJECT_DIR"
pip3 install -r requirements.txt --quiet

# Check if config exists
if [ ! -f "$PROJECT_DIR/config.yaml" ]; then
    echo ""
    echo "Warning: config.yaml not found"
    echo "Please copy config.example.yaml to config.yaml and configure your location"
    echo ""
    echo "  cp config.example.yaml config.yaml"
    echo "  # Edit config.yaml with your coordinates and preferences"
    echo ""
    exit 1
fi

# Create logs directory
mkdir -p "$LOGS_DIR"

# Create launchd directory if it doesn't exist
mkdir -p "$LAUNCHD_DIR"

# Create updater plist (runs daily at 12:05 AM to update schedule)
UPDATER_PLIST="$LAUNCHD_DIR/com.daynightcycle.updater.plist"
sed -e "s|PYTHON_PATH_PLACEHOLDER|$PYTHON_PATH|g" \
    -e "s|PROJECT_DIR_PLACEHOLDER|$PROJECT_DIR|g" \
    "$PROJECT_DIR/scripts/com.daynightcycle.updater.plist" > "$UPDATER_PLIST"

echo "Created updater plist at: $UPDATER_PLIST"

# Run initial schedule update
echo ""
echo "Calculating today's sunrise and sunset times..."
python3 "$PROJECT_DIR/scripts/update_schedule.py"

# Load launchd agents
echo ""
echo "Loading launchd agents..."

# Unload if already loaded (ignore errors)
launchctl unload "$UPDATER_PLIST" 2>/dev/null || true
launchctl unload "$LAUNCHD_DIR/com.daynightcycle.sunrise.plist" 2>/dev/null || true
launchctl unload "$LAUNCHD_DIR/com.daynightcycle.sunset.plist" 2>/dev/null || true

# Load agents
launchctl load "$UPDATER_PLIST"
launchctl load "$LAUNCHD_DIR/com.daynightcycle.sunrise.plist"
launchctl load "$LAUNCHD_DIR/com.daynightcycle.sunset.plist"

echo ""
echo "================================================"
echo "Installation complete!"
echo "================================================"
echo ""
echo "The system will now automatically:"
echo "  • Switch to light mode at sunrise"
echo "  • Switch to dark mode at sunset"
echo "  • Update schedule daily at 12:05 AM"
echo ""
echo "Manual commands:"
echo "  python3 -m day_night_cycle auto     # Apply mode based on current time"
echo "  python3 -m day_night_cycle light    # Force light mode"
echo "  python3 -m day_night_cycle dark     # Force dark mode"
echo "  python3 -m day_night_cycle status   # Show current status"
echo "  python3 -m day_night_cycle next     # Show next transition"
echo ""
echo "To uninstall:"
echo "  ./scripts/uninstall.sh"
echo ""
