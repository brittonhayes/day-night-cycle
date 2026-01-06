#!/bin/bash
set -e

echo "=========================================="
echo "day-night-cycle installation (Go version)"
echo "=========================================="
echo ""

INSTALL_DIR="$HOME/.config/day-night-cycle"
BINARY_NAME="day-night-cycle"
REPO="brittonhayes/day-night-cycle"
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
mkdir -p "$INSTALL_DIR"

# Download binary
echo "Downloading $BINARY_NAME..."
if ! curl -fsSL "$BINARY_URL" -o "$INSTALL_DIR/$BINARY_NAME"; then
    echo "Error: Failed to download binary from $BINARY_URL"
    echo "Please check that the release exists at: https://github.com/$REPO/releases"
    exit 1
fi

chmod +x "$INSTALL_DIR/$BINARY_NAME"
echo "Downloaded to: $INSTALL_DIR/$BINARY_NAME"
echo ""

# Interactive configuration
if [ -f "$INSTALL_DIR/config.yaml" ]; then
    echo "Found existing config.yaml"
    read -p "Do you want to reconfigure? (y/N): " reconfigure
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
    echo "Find your coordinates: https://www.latlong.net/"
    echo "Find your timezone: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones"
    echo ""

    read -p "Enter your latitude (e.g., 46.0645): " latitude
    read -p "Enter your longitude (e.g., -118.3430): " longitude
    read -p "Enter your timezone (e.g., America/Los_Angeles): " timezone

    cat > "$INSTALL_DIR/config.yaml" <<EOF
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

# Generate launchd schedule
echo ""
echo "Generating launchd schedule..."
"$INSTALL_DIR/$BINARY_NAME" --config "$INSTALL_DIR/config.yaml" schedule

# Load launchd agent
PLIST_PATH="$HOME/Library/LaunchAgents/com.daynightcycle.schedule.plist"
echo ""
echo "Loading launchd agent..."
launchctl unload "$PLIST_PATH" 2>/dev/null || true
launchctl load "$PLIST_PATH"

echo ""
echo "==========================================="
echo "Installation complete!"
echo "==========================================="
echo ""
echo "Commands:"
echo "  $INSTALL_DIR/$BINARY_NAME auto    # Apply based on current time"
echo "  $INSTALL_DIR/$BINARY_NAME light   # Force light mode"
echo "  $INSTALL_DIR/$BINARY_NAME dark    # Force dark mode"
echo "  $INSTALL_DIR/$BINARY_NAME status  # Show status"
echo "  $INSTALL_DIR/$BINARY_NAME next    # Show next transition"
echo ""
echo "Optional: Add an alias to your shell config (~/.zshrc or ~/.bashrc):"
echo "  alias dnc='$INSTALL_DIR/$BINARY_NAME --config $INSTALL_DIR/config.yaml'"
echo ""
echo "Configuration file: $INSTALL_DIR/config.yaml"
echo ""
echo "To uninstall:"
echo "  launchctl unload $PLIST_PATH"
echo "  rm -rf $INSTALL_DIR"
echo "  rm $PLIST_PATH"
echo ""
