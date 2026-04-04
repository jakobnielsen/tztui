# tztui

A terminal UI for browsing timezones and managing your system timezone, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- **Dashboard** — shows your favourite timezones with a live clock
- **Search/Browse** — fuzzy-filter all 500+ IANA zones; toggle favourites with `f`
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

### AUR (Arch Linux)

```sh
yay -S tztui-bin
```

### Debian / Ubuntu

Download the `.deb` from the [latest release](https://github.com/jakobnielsen/tztui/releases/latest) and install it:

```sh
sudo dpkg -i tztui_linux_amd64.deb
```

Replace `amd64` with `arm64` if you are on an ARM machine.

### Fedora / RHEL / openSUSE

Download the `.rpm` from the [latest release](https://github.com/jakobnielsen/tztui/releases/latest) and install it:

```sh
# Fedora / RHEL
sudo rpm -i tztui_linux_amd64.rpm

# openSUSE
sudo zypper install tztui_linux_amd64.rpm
```

Replace `amd64` with `arm64` if you are on an ARM machine.

### Homebrew (macOS + Linux)

```sh
brew tap jakobnielsen/tztui
brew install --cask tztui
```

> **macOS Gatekeeper warning**: If macOS blocks the binary with "Apple could not verify 'tztui' is free of malware", run:
> ```sh
> xattr -d com.apple.quarantine $(which tztui)
> ```
> Alternatively, go to **System Settings → Privacy & Security** and click **Allow Anyway**.

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

Changing the system timezone requires elevated privileges. `tztui` will attempt the change directly and, if permission is denied, prompt you for your sudo password inside the TUI — no need to run `sudo tztui`.

- **Linux** — uses `timedatectl set-timezone`
- **macOS** — uses `systemsetup -settimezone`

## Platform support

| Platform | Supported |
|---|---|
| Linux (amd64, arm64) | ✓ |
| macOS (amd64, arm64) | ✓ |
| Windows | Not yet |

## License

MIT
