package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jakobnielsen/tztui/internal/config"
	"github.com/jakobnielsen/tztui/internal/tz"
)

// TickMsg is sent every second to update the clock display.
type TickMsg time.Time

// Dashboard shows a table of favourite timezones with a live clock.
type Dashboard struct {
	cfg      *config.Config
	table    table.Model
	systemTZ string
	width    int
	height   int
	err      string
}

func newDashboard(cfg *config.Config, systemTZ string) Dashboard {
	// Default column widths used before the first WindowSizeMsg arrives.
	cols := defaultDashCols()

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = StyleTableHeader
	s.Selected = StyleTableRowSelected
	t.SetStyles(s)

	d := Dashboard{
		cfg:      cfg,
		table:    t,
		systemTZ: systemTZ,
	}
	d.refreshRows()
	return d
}

func defaultDashCols() []table.Column {
	return []table.Column{
		{Title: "★", Width: 2},
		{Title: "Label", Width: 20},
		{Title: "IANA", Width: 28},
		{Title: "Local time", Width: 21},
		{Title: "Offset", Width: 10},
	}
}

// resizeCols distributes available terminal width across the flexible table columns.
func (d *Dashboard) resizeCols() {
	if d.width == 0 {
		d.table.SetColumns(defaultDashCols())
		return
	}
	// Panel overhead: 2 border + 2 padding = 4. Panel content = d.width-4.
	// Each of 5 table cells has Padding(0,1) = +2 chars each = 10 total overhead.
	// Fixed col widths: star(2) + time(21) + offset(10) = 33.
	// flex = (d.width-4) - 10 - 33 = d.width - 47
	flex := max(20, d.width-47)
	labelW := flex * 40 / 100
	ianaW := flex - labelW
	d.table.SetColumns([]table.Column{
		{Title: "★", Width: 2},
		{Title: "Label", Width: labelW},
		{Title: "IANA", Width: ianaW},
		{Title: "Local time", Width: 21},
		{Title: "Offset", Width: 10},
	})
}

func (d *Dashboard) refreshRows() {
	rows := make([]table.Row, 0, len(d.cfg.Favourites))
	for _, iana := range d.cfg.Favourites {
		t, err := tz.CurrentTime(iana)
		var timeStr, offset, star string
		if err != nil {
			timeStr = "—"
			offset = "—"
		} else {
			timeStr = t.Format("Mon 02 Jan  15:04:05")
			_, off := t.Zone()
			h, m := off/3600, (off%3600)/60
			if m < 0 {
				m = -m
			}
			if m == 0 {
				offset = fmt.Sprintf("UTC%+d", h)
			} else {
				offset = fmt.Sprintf("UTC%+d:%02d", h, m)
			}
		}
		if iana == d.systemTZ {
			star = "●"
		}
		label := tz.LabelForIANA(iana)
		rows = append(rows, table.Row{star, label, iana, timeStr, offset})
	}
	d.table.SetRows(rows)
}

// SelectedIANA returns the IANA identifier of the currently highlighted row,
// or "" if the table is empty.
func (d *Dashboard) SelectedIANA() string {
	row := d.table.SelectedRow()
	if len(row) < 3 {
		return ""
	}
	return row[2]
}

func (d *Dashboard) Init() tea.Cmd {
	return tick()
}

func (d *Dashboard) Update(msg tea.Msg) (Dashboard, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		d.refreshRows()
		return *d, tick()

	case tea.KeyMsg:
		switch {
		case key_matches(msg, DashKeys.RemoveFav):
			iana := d.SelectedIANA()
			if iana != "" {
				d.cfg.RemoveFavourite(iana)
				_ = config.Save(d.cfg)
				d.refreshRows()
			}
			return *d, nil

		case key_matches(msg, DashKeys.SetTZ):
			if iana := d.SelectedIANA(); iana != "" {
				return *d, func() tea.Msg { return SelectMsg{IANA: iana} }
			}
			return *d, nil

		case key_matches(msg, DashKeys.Search):
			// Handled by parent — just pass through.
		}

	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		// Overhead lines: header(2)+tabs(3)+panel borders(2)+title(2)+newline(1)
		// +table header+border(2)+newline(1)+help(2) = 15. Use 16 for safety.
		d.table.SetHeight(max(4, msg.Height-16))
		d.resizeCols()
	}

	var cmd tea.Cmd
	d.table, cmd = d.table.Update(msg)
	return *d, cmd
}

func (d *Dashboard) View() string {
	if len(d.cfg.Favourites) == 0 {
		empty := lipgloss.JoinVertical(lipgloss.Left,
			StylePanelTitle.Render("Favourites"),
			StyleMuted.Render("No favourites yet. Go to Search (tab) and press"),
			StyleMuted.Render(StyleFavStar.Render("f")+" to add a timezone."),
		)
		return StylePanel.Render(empty)
	}

	var sb strings.Builder
	sb.WriteString(StylePanelTitle.Render("Favourites"))
	sb.WriteString("\n")
	sb.WriteString(d.table.View())
	sb.WriteString("\n")
	sb.WriteString(renderHelp("↑/↓", "move", "enter", "set timezone", "f", "remove", "/", "search", "tab", "next view", "●", "current system timezone"))

	return StylePanel.Render(sb.String())
}

// tick returns a command that fires TickMsg after one second.
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// key_matches is a small helper to avoid importing bubbles/key directly in
// every switch statement.
func key_matches(msg tea.KeyMsg, b interface{ Keys() []string }) bool {
	for _, k := range b.Keys() {
		if k == msg.String() {
			return true
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
