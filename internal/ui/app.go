// Package ui contains the Bubble Tea application model and all view models.
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jakobnielsen/tztui/internal/config"
)

type viewID int

const (
	viewDashboard viewID = iota
	viewSearch
	viewPicker
	viewCount // sentinel — keep last
)

func (v viewID) String() string {
	switch v {
	case viewDashboard:
		return "Dashboard"
	case viewSearch:
		return "Search"
	case viewPicker:
		return "Set Timezone"
	}
	return ""
}

// App is the root Bubble Tea model.
type App struct {
	cfg       *config.Config
	systemTZ  string
	version   string
	current   viewID
	dashboard Dashboard
	search    Search
	picker    Picker
	width     int
	height    int
}

// NewApp constructs the root model.
// setter is called first; sudoSetter is called if setter returns ErrPermission.
func NewApp(cfg *config.Config, systemTZ, version string, setter func(string) error, sudoSetter func(string, string) error) App {
	return App{
		cfg:       cfg,
		systemTZ:  systemTZ,
		version:   version,
		current:   viewDashboard,
		dashboard: newDashboard(cfg, systemTZ),
		search:    newSearch(cfg),
		picker:    newPicker(setter, sudoSetter),
	}
}

func (a App) Init() tea.Cmd {
	return tea.Batch(
		a.dashboard.Init(),
		a.search.Init(),
	)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		var c1, c2, c3 tea.Cmd
		a.dashboard, c1 = a.dashboard.Update(msg)
		a.search, c2 = a.search.Update(msg)
		_, c3 = a.picker.Update(msg)
		return a, tea.Batch(c1, c2, c3)

	case TickMsg:
		var cmd tea.Cmd
		a.dashboard, cmd = a.dashboard.Update(msg)
		return a, cmd

	case SelectMsg:
		if msg.IANA == "" {
			// Picker cancelled / dismissed — go back to search.
			if a.current == viewPicker {
				a.current = viewSearch
				a.search.Focus()
			}
			return a, nil
		}
		// User selected a zone from Search — open Picker.
		a.picker.Activate(msg.IANA)
		a.current = viewPicker
		a.search.Blur()
		return a, nil

	case TZChangedMsg:
		a.systemTZ = msg.IANA
		a.dashboard.systemTZ = msg.IANA
		a.dashboard.refreshRows()
		return a, nil

	case tea.KeyMsg:
		// Global quit.
		if key.Matches(msg, Keys.Quit) {
			return a, tea.Quit
		}

		// Navigate tabs (not while picker is open).
		if a.current != viewPicker {
			if key.Matches(msg, Keys.Tab) {
				a.current = (a.current + 1) % (viewCount - 1) // exclude Picker from tab cycle
				a.syncFocus()
				return a, nil
			}
			if key.Matches(msg, Keys.ShiftTab) {
				a.current = ((a.current - 1) + (viewCount - 1)) % (viewCount - 1)
				a.syncFocus()
				return a, nil
			}
		}

		// "/" on dashboard goes to search.
		if a.current == viewDashboard && msg.String() == "/" {
			a.current = viewSearch
			a.syncFocus()
			return a, nil
		}
	}

	// Route to active view.
	var cmd tea.Cmd
	switch a.current {
	case viewDashboard:
		a.dashboard, cmd = a.dashboard.Update(msg)
	case viewSearch:
		a.search, cmd = a.search.Update(msg)
	case viewPicker:
		a.picker, cmd = a.picker.Update(msg)
	}
	return a, cmd
}

func (a *App) syncFocus() {
	if a.current == viewSearch {
		a.search.Focus()
	} else {
		a.search.Blur()
	}
}

func (a App) View() string {
	header := a.renderHeader()
	tabs := a.renderTabs()

	var body string
	switch a.current {
	case viewDashboard:
		body = a.dashboard.View()
	case viewSearch:
		body = a.search.View()
	case viewPicker:
		body = lipgloss.JoinVertical(lipgloss.Left,
			a.search.View(),
			"\n",
			a.picker.View(),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, tabs, body)
}

func (a App) renderHeader() string {
	title := StyleHeader.Render("tztui")
	sysTZ := fmt.Sprintf(" system: %s ", a.systemTZ)
	badge := StyleCurrentTZBadge.Render(sysTZ)
	ver := StyleMuted.Render("  v" + a.version)

	right := lipgloss.JoinHorizontal(lipgloss.Top, badge, ver)

	gap := ""
	usedWidth := lipgloss.Width(title) + lipgloss.Width(right)
	if a.width > usedWidth+2 {
		gap = strings.Repeat(" ", a.width-usedWidth)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, title, gap, right)
}

func (a App) renderTabs() string {
	tabs := make([]string, 0, 2)
	for _, v := range []viewID{viewDashboard, viewSearch} {
		label := v.String()
		if v == a.current || (a.current == viewPicker && v == viewSearch) {
			tabs = append(tabs, StyleTabActive.Render(label))
		} else {
			tabs = append(tabs, StyleTabInactive.Render(label))
		}
	}
	return StyleTabBar.Render(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
}
