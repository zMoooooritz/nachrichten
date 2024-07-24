package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/nachrichten/pkg/config"
)

type HelpState int

const (
	HS_HIDDEN HelpState = iota
	HS_NORMAL
	HS_ALL
)

type Helper struct {
	model  help.Model
	keymap KeyMap
	state  HelpState
}

func NewHelper(s config.Style, keyMap KeyMap, state HelpState) *Helper {
	h := Helper{
		model:  help.New(),
		keymap: keyMap,
		state:  state,
	}
	h.model.FullSeparator = " • "
	h.model.Styles.ShortKey = s.InactiveStyle
	h.model.Styles.FullKey = s.InactiveStyle

	return &h
}

func (h *Helper) View() string {
	if h.IsVisible() {
		return "\n" + lipgloss.NewStyle().Width(h.model.Width).AlignHorizontal(lipgloss.Center).Render(h.model.View(h.keymap))
	}
	return ""
}

func (h *Helper) NextState() {
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
