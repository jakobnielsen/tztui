package ui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jakobnielsen/tztui/internal/tz"
)

// TZChangedMsg is emitted after a successful system timezone change.
type TZChangedMsg struct {
	IANA string
}

// pickerState tracks what stage of the confirmation flow we are in.
type pickerState int

const (
	pickerConfirm pickerState = iota
	pickerWorking
	pickerPassword     // sudo password needed
	pickerPasswordBusy // waiting for sudo to complete
	pickerDone
	pickerError
)

// Picker is a confirmation dialog for changing the system timezone.
type Picker struct {
	iana       string
	label      string
	state      pickerState
	message    string
	setter     func(string) error
	sudoSetter func(string, string) error
	passInput  textinput.Model
}

func newPicker(setter func(string) error, sudoSetter func(string, string) error) Picker {
	ti := textinput.New()
	ti.Placeholder = "sudo password"
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	ti.CharLimit = 256
	ti.Width = 30
	return Picker{setter: setter, sudoSetter: sudoSetter, passInput: ti}
}

// Activate sets the zone to confirm and resets state.
func (p *Picker) Activate(iana string) {
	p.iana = iana
	p.label = tz.LabelForIANA(iana)
	p.state = pickerConfirm
	p.message = ""
	p.passInput.SetValue("")
	p.passInput.Blur()
}

func (p *Picker) Init() tea.Cmd { return nil }

func (p *Picker) Update(msg tea.Msg) (Picker, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch p.state {
		case pickerConfirm:
			switch {
			case key_matches(msg, PickKeys.Confirm):
				p.state = pickerWorking
				return *p, p.applyChange()
			case key_matches(msg, PickKeys.Cancel):
				return *p, func() tea.Msg { return SelectMsg{IANA: ""} }
			}
		case pickerPassword:
			switch msg.String() {
			case "enter":
				p.state = pickerPasswordBusy
				p.passInput.Blur()
				return *p, p.applyChangeWithSudo()
			case "esc":
				return *p, func() tea.Msg { return SelectMsg{IANA: ""} }
			}
		case pickerDone, pickerError:
			return *p, func() tea.Msg { return SelectMsg{IANA: ""} }
		}

	case tzSetResultMsg:
		if msg.err != nil {
			if errors.Is(msg.err, tz.ErrPermission) {
				// First attempt failed due to insufficient privileges — ask for password.
				p.state = pickerPassword
				p.passInput.Focus()
				return *p, textinput.Blink
			}
			p.state = pickerError
			p.message = msg.err.Error()
		} else {
			p.state = pickerDone
			p.message = fmt.Sprintf("System timezone changed to %s", p.iana)
			return *p, func() tea.Msg { return TZChangedMsg{IANA: p.iana} }
		}
	}

	// Forward events to password input when it is active.
	if p.state == pickerPassword {
		var cmd tea.Cmd
		p.passInput, cmd = p.passInput.Update(msg)
		return *p, cmd
	}

	return *p, nil
}

type tzSetResultMsg struct{ err error }

func (p *Picker) applyChange() tea.Cmd {
	iana := p.iana
	setter := p.setter
	return func() tea.Msg {
		err := setter(iana)
		return tzSetResultMsg{err: err}
	}
}

func (p *Picker) applyChangeWithSudo() tea.Cmd {
	iana := p.iana
	password := p.passInput.Value()
	sudoSetter := p.sudoSetter
	return func() tea.Msg {
		err := sudoSetter(iana, password)
		return tzSetResultMsg{err: err}
	}
}

func (p *Picker) View() string {
	t, _ := tz.CurrentTime(p.iana)
	timeStr := t.Format("Mon 02 Jan 2006  15:04:05 MST")

	var body strings.Builder

	switch p.state {
	case pickerConfirm:
		body.WriteString(StylePanelTitle.Render("Set System Timezone"))
		body.WriteString("\n\n")
		body.WriteString(fmt.Sprintf("%-14s %s\n", "Timezone:", StyleClock.Render(p.label)))
		body.WriteString(fmt.Sprintf("%-14s %s\n", "IANA ID:", StyleMuted.Render(p.iana)))
		body.WriteString(fmt.Sprintf("%-14s %s\n", "Current time:", StyleClock.Render(timeStr)))
		body.WriteString("\n")
		body.WriteString(StyleSuccess.Render("Press enter / y to confirm"))
		body.WriteString("\n")
		body.WriteString(StyleError.Render("Press esc / n to cancel"))

	case pickerWorking, pickerPasswordBusy:
		body.WriteString(StyleMuted.Render("Applying timezone changeâ¦"))

	case pickerPassword:
		body.WriteString(StylePanelTitle.Render("Sudo password required"))
		body.WriteString("\n\n")
		body.WriteString(fmt.Sprintf("%-14s %s\n", "Timezone:", StyleClock.Render(p.label)))
		body.WriteString("\n")
		body.WriteString(StyleMuted.Render("Enter your sudo password to apply the change:"))
		body.WriteString("\n")
		body.WriteString(StyleInput.Render(p.passInput.View()))
		body.WriteString("\n\n")
		body.WriteString(StyleHelp.Render("enter confirm  esc cancel"))

	case pickerDone:
		body.WriteString(StyleSuccess.Render("â " + p.message))
		body.WriteString("\n\n")
		body.WriteString(StyleHelp.Render("Press any key to continue."))

	case pickerError:
		body.WriteString(StyleError.Render("â Failed to change timezone"))
		body.WriteString("\n\n")
		body.WriteString(p.message)
		body.WriteString("\n\n")
		body.WriteString(StyleHelp.Render("Press any key to go back."))
	}

	return StyleConfirmBox.Render(body.String())
}
