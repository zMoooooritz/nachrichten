package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/zMoooooritz/nachrichten/pkg/config"
)

func NewDotSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return s
}

func NewHelper(s config.Style) help.Model {
	h := help.New()
	h.FullSeparator = " • "
	h.Styles.ShortKey = s.InactiveStyle
	h.Styles.FullKey = s.InactiveStyle
	return h
}
