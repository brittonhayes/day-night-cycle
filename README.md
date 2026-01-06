# day-night-cycle

Switch application themes at sunrise and sunset.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/brittonhayes/day-night-cycle/main/install-go.sh | bash
```

Or:

```bash
go install github.com/brittonhayes/day-night-cycle@latest
```

Or download a binary from the [releases page](https://github.com/brittonhayes/day-night-cycle/releases).

## Configure

Edit `~/.config/day-night-cycle/config.yaml`:

```yaml
location:
  latitude: 46.0645
  longitude: -118.3430
  timezone: "America/Los_Angeles"

plugins:
  - name: iterm2
    enabled: true
    light_preset: "Light Background"
    dark_preset: "Dark Background"

  - name: cursor
    enabled: true
    light_theme: "GitHub Light"
    dark_theme: "Dark Modern"

  - name: claude-code
    enabled: true

  - name: neovim
    enabled: true
    light_colorscheme: "github_light"
    dark_colorscheme: "github_dark_default"

  - name: macos-system
    enabled: false
```

## Use

```bash
day-night-cycle auto      # apply mode for current time
day-night-cycle light     # force light mode
day-night-cycle dark      # force dark mode
day-night-cycle status    # show current status
day-night-cycle next      # show next transition
day-night-cycle schedule  # generate launchd schedule
```

## Build

```bash
make build
```
