#!/bin/bash
# Uninstallation script for Day/Night Cycle automation

LAUNCHD_DIR="$HOME/Library/LaunchAgents"

echo "Uninstalling Day/Night Cycle automation..."
echo ""

# Unload and remove launchd agents
echo "Removing launchd agents..."

launchctl unload "$LAUNCHD_DIR/com.daynightcycle.updater.plist" 2>/dev/null || true
launchctl unload "$LAUNCHD_DIR/com.daynightcycle.sunrise.plist" 2>/dev/null || true
launchctl unload "$LAUNCHD_DIR/com.daynightcycle.sunset.plist" 2>/dev/null || true

rm -f "$LAUNCHD_DIR/com.daynightcycle.updater.plist"
rm -f "$LAUNCHD_DIR/com.daynightcycle.sunrise.plist"
rm -f "$LAUNCHD_DIR/com.daynightcycle.sunset.plist"

echo ""
echo "Uninstallation complete!"
echo ""
echo "Note: Python dependencies and project files are still present"
echo "Remove the project directory manually if desired"
echo ""
