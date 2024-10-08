package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
)

type HomeSelector struct {
	BaseSelector
}

func NewHomeSelector(selector BaseSelector) *HomeSelector {
	return &HomeSelector{
		BaseSelector: selector,
	}
}

func (s HomeSelector) Init() tea.Cmd {
	return nil
}

func (s *HomeSelector) Update(msg tea.Msg) (Selector, tea.Cmd) {
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
	return &HomeSelector{BaseSelector: bs}, tea.Batch(cmds...)
}

func (s HomeSelector) View() string {
	s.list.SetSize(s.width, s.height)

	return s.list.View()
}
