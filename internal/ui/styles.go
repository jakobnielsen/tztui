package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	colorPrimary   = lipgloss.Color("#7C3AED")
	colorSecondary = lipgloss.Color("#6D28D9")
	colorAccent    = lipgloss.Color("#A78BFA")
	colorMuted     = lipgloss.Color("#6B7280")
	colorText      = lipgloss.Color("#F9FAFB")
	colorTextDim   = lipgloss.Color("#9CA3AF")
	colorBg        = lipgloss.Color("#111827")
	colorSuccess   = lipgloss.Color("#10B981")
	colorError     = lipgloss.Color("#EF4444")
	colorFav       = lipgloss.Color("#F59E0B")

	StyleApp = lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorText)

	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent).
			PaddingBottom(1)

	StyleTabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(colorPrimary)

	StyleTabInactive = lipgloss.NewStyle().
				Foreground(colorMuted).
				PaddingLeft(1).
				PaddingRight(1)

	StyleTabBar = lipgloss.NewStyle().
			PaddingBottom(1)

	StylePanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSecondary).
			Padding(0, 1)

	StylePanelTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent).
			PaddingBottom(1)

	StyleTableHeader = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorAccent).
				Padding(0, 1).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(colorMuted)

	StyleTableRow = lipgloss.NewStyle().
			Foreground(colorText)

	StyleTableRowSelected = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorBg).
				Background(colorPrimary)

	StyleClock = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText)

	StyleFavStar = lipgloss.NewStyle().
			Foreground(colorFav)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true)

	StyleError = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true)

	StyleMuted = lipgloss.NewStyle().
			Foreground(colorMuted)

	StyleHelp = lipgloss.NewStyle().
			Foreground(colorTextDim).
			PaddingTop(1)

	StyleInput = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1)

	StyleConfirmBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(1, 3).
			Width(50)

	StyleCurrentTZBadge = lipgloss.NewStyle().
				Foreground(colorBg).
				Background(colorAccent).
				Bold(true).
				Padding(0, 1)

	styleHelpKey = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)
)

// renderHelp formats a help string where keys and descriptions alternate,
// separated by spaces. Keys are rendered in accent colour, descriptions in dim.
// Pass pairs of (key, description), e.g. renderHelp("↑/↓","move", "f","favourite").
func renderHelp(pairs ...string) string {
	parts := make([]string, 0, len(pairs))
	for i, s := range pairs {
		if i%2 == 0 {
			parts = append(parts, styleHelpKey.Render(s))
		} else {
			parts = append(parts, StyleMuted.Render(s))
		}
	}
	return StyleHelp.Render(strings.Join(parts, " "))
}
