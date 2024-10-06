package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HelpState int

const (
	HS_HIDDEN HelpState = iota
	HS_NORMAL
	HS_ALL
)

type Helper struct {
	shared *SharedState
	model  help.Model
	state  HelpState
}

func NewHelper(shared *SharedState, hstate HelpState) *Helper {
	h := Helper{
		shared: shared,
		model:  help.New(),
		state:  hstate,
	}
	h.model.FullSeparator = " • "
	h.model.Styles.ShortKey = shared.style.InactiveStyle
	h.model.Styles.FullKey = shared.style.InactiveStyle

	return &h
}

func (h Helper) View() string {
	if h.IsVisible() {
		return "\n" + lipgloss.NewStyle().Width(h.model.Width).AlignHorizontal(lipgloss.Center).Render(h.model.View(h.shared.keymap))
	}
	return ""
}

func (h *Helper) Update(msg tea.Msg) (*Helper, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if h.shared.mode == INSERT_MODE {
			break
		}

		switch {
		case key.Matches(msg, h.shared.keymap.help):
			h.nextState()
		}
	}
	return h, nil
}

func (h *Helper) nextState() {
	h.state = (h.state + 1) % 3
	if h.state == HS_NORMAL {
		h.model.ShowAll = false
	}
	if h.state == HS_ALL {
		h.model.ShowAll = true
	}
}

func (h *Helper) Height() int {
	if h.IsVisible() {
		return 2
	}
	return 0
}

func (h *Helper) SetWidth(width int) {
	h.model.Width = width
}

func (h *Helper) IsVisible() bool {
	return h.state != HS_HIDDEN
}
