package main

import (
	_ "time/tzdata"

	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jakobnielsen/tztui/internal/config"
	"github.com/jakobnielsen/tztui/internal/tz"
	"github.com/jakobnielsen/tztui/internal/ui"
)

// version is set at build time via -ldflags "-X main.version=<tag>".
var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println("tztui", version)
		os.Exit(0)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "tztui: load config: %v\n", err)
		os.Exit(1)
	}

	systemTZ, err := tz.GetSystemTZ()
	if err != nil {
		// Non-fatal: fall back to UTC so the TUI still works.
		systemTZ = "UTC"
	}

	app := ui.NewApp(cfg, systemTZ, version, tz.SetSystemTZ, tz.SetSystemTZWithSudo)

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "tztui: %v\n", err)
		os.Exit(1)
	}
}
