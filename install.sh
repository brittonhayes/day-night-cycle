#!/bin/bash
set -e

CONFIG_DIR="$HOME/.config/day-night-cycle"
BINARY_INSTALL_DIR="/usr/local/bin"
BINARY_NAME="day-night-cycle"
REPO="brittonhayes/day-night-cycle"
PLIST_PATH="$HOME/Library/LaunchAgents/com.daynightcycle.schedule.plist"

# Handle uninstall
if [ "$1" = "--uninstall" ]; then
    echo "=========================================="
    echo "day-night-cycle uninstallation"
    echo "=========================================="
    echo ""

    # Unload launchd agent
    if [ -f "$PLIST_PATH" ]; then
        echo "Unloading launchd agent..."
        launchctl unload "$PLIST_PATH" 2>/dev/null || true
        rm "$PLIST_PATH"
        echo "Removed: $PLIST_PATH"
    fi

    # Remove binary
    if [ -f "$BINARY_INSTALL_DIR/$BINARY_NAME" ]; then
        echo "Removing binary..."
        rm "$BINARY_INSTALL_DIR/$BINARY_NAME"
        echo "Removed: $BINARY_INSTALL_DIR/$BINARY_NAME"
    fi

    # Remove configuration directory
    if [ -d "$CONFIG_DIR" ]; then
        echo "Removing configuration directory..."
        rm -rf "$CONFIG_DIR"
        echo "Removed: ~/.config/day-night-cycle"
    fi

    echo ""
    echo "=========================================="
    echo "Uninstallation complete!"
    echo "=========================================="
    exit 0
fi

echo "=========================================="
echo "day-night-cycle installation (Go version)"
echo "=========================================="
echo ""

ARCH=$(uname -m)

# Detect architecture
if [ "$ARCH" = "arm64" ]; then
    BINARY_SUFFIX="darwin-arm64"
    echo "Detected: Apple Silicon (arm64)"
elif [ "$ARCH" = "x86_64" ]; then
    BINARY_SUFFIX="darwin-amd64"
    echo "Detected: Intel (amd64)"
else
    echo "Error: Unsupported architecture: $ARCH"
    exit 1
fi

# Get latest release version
echo "Fetching latest release..."
LATEST_VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "Error: Could not determine latest version"
    echo "You can manually download from: https://github.com/$REPO/releases"
    exit 1
fi

echo "Latest version: $LATEST_VERSION"
echo ""

# Build download URL
BINARY_URL="https://github.com/$REPO/releases/download/${LATEST_VERSION}/${BINARY_NAME}-${BINARY_SUFFIX}"

# Create config directory
mkdir -p "$CONFIG_DIR"

# Download binary
echo "Downloading $BINARY_NAME..."
TEMP_BINARY="/tmp/$BINARY_NAME"
if ! curl -fsSL "$BINARY_URL" -o "$TEMP_BINARY"; then
    echo "Error: Failed to download binary from $BINARY_URL"
    echo "Please check that the release exists at: https://github.com/$REPO/releases"
    exit 1
fi

# Install binary to /usr/local/bin
chmod +x "$TEMP_BINARY"
mv "$TEMP_BINARY" "$BINARY_INSTALL_DIR/$BINARY_NAME"
echo "Installed to: $BINARY_INSTALL_DIR/$BINARY_NAME"
echo ""

# Interactive configuration
if [ -f "$CONFIG_DIR/config.yaml" ]; then
    echo "Found existing config.yaml"
    read -p "Do you want to reconfigure? (y/N): " reconfigure </dev/tty
    if [[ ! $reconfigure =~ ^[Yy]$ ]]; then
        echo "Using existing configuration"
        SKIP_CONFIG=1
    fi
fi

if [ -z "$SKIP_CONFIG" ]; then
    echo "==========================================="
    echo "Configuration"
    echo "==========================================="
    echo ""
    echo "Find your coordinates: https://www.lat-long-coordinates.com/"
    echo "Find your timezone: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones"
    echo ""

    # Prompt for required configuration with validation
    # Use /dev/tty to read from terminal when script is piped
    while [ -z "$latitude" ]; do
        read -p "Enter your latitude (e.g., 46.0645): " latitude </dev/tty
        if [ -z "$latitude" ]; then
            echo "Error: Latitude is required"
        fi
    done

    while [ -z "$longitude" ]; do
        read -p "Enter your longitude (e.g., -118.3430): " longitude </dev/tty
        if [ -z "$longitude" ]; then
            echo "Error: Longitude is required"
        fi
    done

    while [ -z "$timezone" ]; do
        read -p "Enter your timezone (e.g., America/Los_Angeles): " timezone </dev/tty
        if [ -z "$timezone" ]; then
            echo "Error: Timezone is required"
        fi
    done

    cat > "$CONFIG_DIR/config.yaml" <<EOF
# yaml-language-server: \$schema=https://raw.githubusercontent.com/brittonhayes/day-night-cycle/main/config.schema.json
location:
  latitude: $latitude
  longitude: $longitude
  timezone: "$timezone"
  # Optional: Adjust when transitions occur (negative = earlier, positive = later)
  # dayOffset: "30m"      # Start day mode 30min after sunrise
  # nightOffset: "-1h"    # Start night mode 1 hour before sunset

plugins:
  - name: macos-system
    enabled: true

  # Uncomment and configure plugins as needed:
  # - name: iterm2
  #   enabled: true
  #   day: "Light Background"
  #   night: "Dark Background"

  # - name: claude-code
  #   enabled: true

  # - name: cursor
  #   enabled: true
  #   day: "Light Modern"
  #   night: "Cursor Dark"

  # - name: neovim
  #   enabled: true

  # - name: sublime
  #   enabled: true
  #   day: "Breakers"
  #   night: "Mariana"

  # - name: pycharm
  #   enabled: true
  #   day: "IntelliJ Light"
  #   night: "Darcula"
EOF

    echo ""
    echo "Configuration saved to ~/.config/day-night-cycle/config.yaml"
    echo "You can edit this file later to customize plugin settings"
fi

# Generate launchd schedule
echo ""
echo "Generating launchd schedule..."
if ! "$BINARY_NAME" --config "$CONFIG_DIR/config.yaml" schedule; then
    echo ""
    echo "Error: Failed to generate launchd schedule"
    echo "Please check your configuration file at: ~/.config/day-night-cycle/config.yaml"
    echo "Make sure all values are properly set (latitude, longitude, timezone)"
    echo ""
    echo "You can manually edit the config and run:"
    echo "  $BINARY_NAME --config ~/.config/day-night-cycle/config.yaml schedule"
    exit 1
fi

# Load launchd agent
echo ""
echo "Setting up automatic scheduling..."
launchctl unload "$PLIST_PATH" 2>/dev/null || true
launchctl load "$PLIST_PATH"
echo "Automatic theme switching enabled"

echo ""
echo "==========================================="
echo "Installation complete!"
echo "==========================================="
echo ""
echo "Commands:"
echo "  $BINARY_NAME auto    # Apply based on current time"
echo "  $BINARY_NAME light   # Force light mode"
echo "  $BINARY_NAME dark    # Force dark mode"
echo "  $BINARY_NAME status  # Show status"
echo "  $BINARY_NAME next    # Show next transition"
echo ""
echo "Binary location: $BINARY_INSTALL_DIR/$BINARY_NAME"
echo "Configuration file: ~/.config/day-night-cycle/config.yaml"
echo ""
echo "To uninstall:"
echo "  curl -fsSL https://raw.githubusercontent.com/$REPO/main/install.sh | bash -s -- --uninstall"
echo ""
