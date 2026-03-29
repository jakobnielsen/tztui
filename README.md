# tztui

A terminal UI for browsing timezones and managing your system timezone, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- **Dashboard** — shows your favourite timezones with a live clock
- **Search/Browse** — fuzzy-filter all ~120 curated IANA zones; toggle favourites with `f`
- **Set Timezone** — confirm and apply a new system timezone (requires root / sudo)
- Self-contained binary — embeds the IANA timezone database via `time/tzdata`

## Keybindings

| Key | Action |
|---|---|
| `tab` / `shift+tab` | Switch between Dashboard and Search views |
| `↑` / `↓` (or `k` / `j`) | Move selection |
| `f` | Toggle favourite (Search) / Remove favourite (Dashboard) |
| `enter` | Open timezone picker |
| `enter` / `y` | Confirm timezone change (Picker) |
| `esc` / `n` | Cancel / go back |
| `/` | Jump to Search from Dashboard |
| `q` / `ctrl+c` | Quit |

## Installation

### Homebrew (macOS + Linux)

```sh
brew tap jakobnielsen/tztui
brew install tztui
```

### AUR (Arch Linux)

```sh
yay -S tztui
```

### go install

```sh
go install github.com/jakobnielsen/tztui/cmd/tztui@latest
```

## Building from source

Requires Go 1.22+.

```sh
git clone https://github.com/jakobnielsen/tztui
cd tztui
make build
./tztui
```

## Configuration

Favourites are persisted to:

```
~/.config/tztui/config.toml
```

Example:

```toml
favourites = ["Europe/London", "Asia/Tokyo", "America/New_York"]
```

## Changing the system timezone

Requires elevated privileges:

- **Linux** — uses `timedatectl set-timezone`; run `tztui` with `sudo` or ensure your user is in the `wheel` group with `sudo` rights
- **macOS** — uses `systemsetup -settimezone`; run `sudo tztui`

## Platform support

| Platform | Supported |
|---|---|
| Linux (amd64, arm64) | ✓ |
| macOS (amd64, arm64) | ✓ |
| Windows | Not yet |

## License

MIT
