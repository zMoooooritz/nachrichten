package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
)

type NewsSelector struct {
	BaseSelector
}

func NewNewsSelector(selector BaseSelector) *NewsSelector {
	return &NewsSelector{
		BaseSelector: selector,
	}
}

func (s NewsSelector) Init() tea.Cmd {
	return nil
}

func (s *NewsSelector) Update(msg tea.Msg) (Selector, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tagesschau.News:
		news := tagesschau.News(msg)
		if s.selectorType == ST_NATIONAL {
			s.articles = news.NationalNews
		} else if s.selectorType == ST_REGIONAL {
			s.articles = news.RegionalNews
		}
		s.rebuildList()
	case tea.KeyMsg:
		if s.isFocused && s.isVisible {
			s.list, cmd = s.list.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	bs, cmd := s.BaseSelector.Update(msg)
	cmds = append(cmds, cmd)
	return &NewsSelector{BaseSelector: bs}, tea.Batch(cmds...)
}

func (s NewsSelector) View() string {
	s.list.SetSize(s.width, s.height)

	return s.list.View()
}
