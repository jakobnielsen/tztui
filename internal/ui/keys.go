package ui

import "github.com/charmbracelet/bubbles/key"

// GlobalKeys are active in every view.
type GlobalKeys struct {
	Quit     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Help     key.Binding
}

var Keys = GlobalKeys{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c / q", "quit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next view"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev view"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

// DashboardKeys are additional bindings active on the Dashboard view.
type DashboardKeys struct {
	RemoveFav key.Binding
	SetTZ     key.Binding
	Search    key.Binding
}

var DashKeys = DashboardKeys{
	RemoveFav: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "remove favourite"),
	),
	SetTZ: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "set as system timezone"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "go to search"),
	),
}

// SearchKeys are additional bindings active on the Search view.
type SearchKeys struct {
	ToggleFav   key.Binding
	Pick        key.Binding
	ClearSearch key.Binding
}

var SearchKeys_ = SearchKeys{
	ToggleFav: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "toggle favourite"),
	),
	Pick: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "set as system timezone"),
	),
	ClearSearch: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear / back"),
	),
}

// PickerKeys are bindings active on the Picker confirmation dialog.
type PickerKeys struct {
	Confirm key.Binding
	Cancel  key.Binding
}

var PickKeys = PickerKeys{
	Confirm: key.NewBinding(
		key.WithKeys("enter", "y"),
		key.WithHelp("enter / y", "confirm"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc", "n"),
		key.WithHelp("esc / n", "cancel"),
	),
}
