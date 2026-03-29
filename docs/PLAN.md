# tztui — Implementation Plan

A Go TUI application using Bubble Tea + Bubbles + Lip Gloss for browsing and managing system timezones.

## Key decisions

| Decision | Choice |
|---|---|
| TUI framework | Bubble Tea (charmbracelet/bubbletea) |
| Styling | Lip Gloss (charmbracelet/lipgloss) |
| Config format | TOML |
| Config location | `$XDG_CONFIG_HOME/tztui/config.toml` (default `~/.config/tztui/`) |
| Timezone DB | stdlib `time/tzdata` embedded — no system dependency |
| System TZ change (Linux) | `timedatectl set-timezone <iana>` |
| System TZ change (macOS) | `systemsetup -settimezone <iana>` |
| Windows | Out of scope (build tags make future addition straightforward) |
| Packaging | GoReleaser → Homebrew tap + AUR PKGBUILD |

## Directory structure

```
tztui/
├── cmd/tztui/
│   └── main.go                 # entry point, --version flag
├── internal/
│   ├── config/
│   │   ├── config.go           # Load/Save/AddFavourite/RemoveFavourite
│   │   └── config_test.go
│   ├── tz/
│   │   ├── zones.go            # curated IANA zone list, CurrentTime, LabelForIANA
│   │   ├── system_linux.go     # (//go:build linux) GetSystemTZ / SetSystemTZ
│   │   ├── system_darwin.go    # (//go:build darwin) GetSystemTZ / SetSystemTZ
│   │   └── tz_test.go
│   └── ui/
│       ├── app.go              # root model, tab routing, window resize
│       ├── dashboard.go        # favourites table with live clock (tea.Tick)
│       ├── search.go           # filterable zone list, toggle favourites
│       ├── picker.go           # confirmation dialog, calls tz.SetSystemTZ
│       ├── styles.go           # all Lip Gloss style constants
│       └── keys.go             # key.Binding definitions per view
├── docs/
│   └── PLAN.md                 # this file
├── Formula/
│   └── tztui.rb                # Homebrew formula (sha256 filled by GoReleaser)
├── .goreleaser.yaml            # multi-platform release config
├── PKGBUILD                    # Arch Linux AUR package
├── Makefile                    # build / test / lint / install targets
├── go.mod
└── README.md
```

## Architecture overview

### Model hierarchy (Bubble Tea)

```
App  (root model, ui/app.go)
├── Dashboard   (viewDashboard)
├── Search      (viewSearch)
└── Picker      (viewPicker — overlays Search)
```

`App.Update` routes messages:
- `TickMsg` → Dashboard (clock refresh)
- `SelectMsg{IANA}` → activates Picker, switches view
- `TZChangedMsg` → updates `App.systemTZ` + Dashboard badge
- `tea.KeyMsg` → global quit/tab, then forwarded to active view

### Platform adapter interface

Both `system_linux.go` and `system_darwin.go` expose:

```go
func GetSystemTZ() (string, error)
func SetSystemTZ(iana string) error
```

`SetSystemTZ` is passed as a `func(string) error` into `ui.NewApp`, enabling injection of a mock in tests without requiring sudo.

### Config schema (`config.toml`)

```toml
favourites = ["Europe/London", "Asia/Tokyo", "America/New_York"]
```

## Packaging

### Homebrew (macOS + Linux)

```sh
brew tap jakobnielsen/tztui
brew install tztui
```

The tap repository `jakobnielsen/homebrew-tztui` is automatically updated by GoReleaser on each release (requires `HOMEBREW_TAP_GITHUB_TOKEN` secret).

### AUR (Arch Linux)

```sh
yay -S tztui
# or
paru -S tztui
```

The `PKGBUILD` at the repo root is also published under `tztui` in the AUR.

### Manual / go install

```sh
go install github.com/jakobnielsen/tztui/cmd/tztui@latest
```

## Verification checklist

- [ ] `make test` — all unit tests pass (no sudo required)
- [ ] `make build` — binary produced; `./tztui` launches in alt-screen
- [ ] Add a favourite in Search → verify `~/.config/tztui/config.toml` is updated
- [ ] `./tztui` dashboard shows live clock updating every second
- [ ] Picker → confirm → system timezone changes (`timedatectl` / `systemsetup`)
- [ ] `goreleaser build --snapshot --clean` builds all four target archives
- [ ] `brew install --build-from-source Formula/tztui.rb` succeeds
- [ ] `tztui --version` prints the correct version string

## Open TODO / future ideas

- [ ] `.goreleaser.yaml` `sha256sums` in `Formula/tztui.rb` are templated — do not edit manually
- [ ] Windows: add `system_windows.go` using `tzutil /s <iana>` and `(//go:build windows)` tag
- [ ] Configurable colour theme via config.toml
- [ ] `--list` flag to print all zones as JSON/TSV for scripting
