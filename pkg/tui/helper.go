package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/zMoooooritz/nachrichten/pkg/config"
)

func EmptyLists(s config.Style, count int) []list.Model {
	var lists []list.Model
	for i := 0; i < count; i++ {
		newList := list.New([]list.Item{}, NewNewsDelegate(s), 0, 0)
		newList.SetFilteringEnabled(false)
		newList.SetShowTitle(true)
		newList.SetShowStatusBar(false)
		newList.SetShowHelp(false)
		lists = append(lists, newList)
	}
	return lists
}

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
