package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jakobnielsen/tztui/internal/config"
	"github.com/jakobnielsen/tztui/internal/tz"
)

// SelectMsg is emitted when the user presses Enter on a zone in Search.
// The parent model routes this to the Picker.
type SelectMsg struct {
	IANA string
}

// Search provides a filterable list of all known IANA timezones.
type Search struct {
	cfg       *config.Config
	input     textinput.Model
	all       []tz.Zone
	filtered  []tz.Zone
	cursor    int
	width     int
	height    int
	visHeight int
}

func newSearch(cfg *config.Config) Search {
	ti := textinput.New()
	ti.Placeholder = "Filter timezones…"
	ti.CharLimit = 60
	ti.Width = 40

	s := Search{
		cfg:       cfg,
		input:     ti,
		all:       tz.ListZones(),
		visHeight: 15,
	}
	s.applyFilter()
	return s
}

func (s *Search) applyFilter() {
	q := strings.ToLower(strings.TrimSpace(s.input.Value()))
	if q == "" {
		s.filtered = make([]tz.Zone, len(s.all))
		copy(s.filtered, s.all)
	} else {
		s.filtered = s.filtered[:0]
		for _, z := range s.all {
			if strings.Contains(strings.ToLower(z.IANA), q) ||
				strings.Contains(strings.ToLower(z.Label), q) {
				s.filtered = append(s.filtered, z)
			}
		}
	}
	// Clamp cursor.
	if s.cursor >= len(s.filtered) {
		s.cursor = max(0, len(s.filtered)-1)
	}
}

func (s *Search) SelectedZone() *tz.Zone {
	if len(s.filtered) == 0 {
		return nil
	}
	z := s.filtered[s.cursor]
	return &z
}

func (s *Search) Focus() {
	s.input.Focus()
}

func (s *Search) Blur() {
	s.input.Blur()
}

func (s *Search) Init() tea.Cmd {
	return textinput.Blink
}

func (s *Search) Update(msg tea.Msg) (Search, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
			return *s, nil

		case "down", "j":
			if s.cursor < len(s.filtered)-1 {
				s.cursor++
			}
			return *s, nil

		case "f":
			if z := s.SelectedZone(); z != nil {
				if s.cfg.IsFavourite(z.IANA) {
					s.cfg.RemoveFavourite(z.IANA)
				} else {
					s.cfg.AddFavourite(z.IANA)
				}
				_ = config.Save(s.cfg)
			}
			return *s, nil

		case "enter":
			if z := s.SelectedZone(); z != nil {
				return *s, func() tea.Msg { return SelectMsg{IANA: z.IANA} }
			}
			return *s, nil

		case "esc":
			s.input.SetValue("")
			s.applyFilter()
			return *s, nil
		}

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.visHeight = max(5, msg.Height-14)
	}

	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	s.applyFilter()
	return *s, cmd
}

func (s *Search) View() string {
	var sb strings.Builder
	sb.WriteString(StylePanelTitle.Render("Search & Browse"))
	sb.WriteString("\n")
	sb.WriteString(StyleInput.Render(s.input.View()))
	sb.WriteString("\n\n")

	if len(s.filtered) == 0 {
		sb.WriteString(StyleMuted.Render("No results."))
		return StylePanel.Render(sb.String())
	}

	// Windowing so we don't render thousands of lines.
	start := 0
	end := len(s.filtered)
	if end > s.visHeight {
		half := s.visHeight / 2
		start = s.cursor - half
		if start < 0 {
			start = 0
		}
		end = start + s.visHeight
		if end > len(s.filtered) {
			end = len(s.filtered)
			start = max(0, end-s.visHeight)
		}
	}

	for i := start; i < end; i++ {
		z := s.filtered[i]
		fav := "  "
		if s.cfg.IsFavourite(z.IANA) {
			fav = StyleFavStar.Render("★ ")
		}
		t, _ := tz.CurrentTime(z.IANA)
		timeStr := t.Format("15:04:05")

		line := lipgloss.JoinHorizontal(lipgloss.Top,
			fav,
			padRight(z.Label, 22),
			StyleMuted.Render(padRight(z.IANA, 30)),
			StyleClock.Render(timeStr),
		)

		if i == s.cursor {
			line = StyleTableRowSelected.Render(line)
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	if len(s.filtered) > s.visHeight {
		sb.WriteString(StyleMuted.Render(
			strings.Repeat("─", 40),
		))
		sb.WriteString("\n")
	}

	sb.WriteString(StyleHelp.Render("↑/↓ move  f toggle favourite  enter set system TZ  esc clear"))

	return StylePanel.Render(sb.String())
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s[:n]
	}
	return s + strings.Repeat(" ", n-len(s))
}
