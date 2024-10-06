package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zMoooooritz/nachrichten/pkg/config"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
)

type SelectorType int

const (
	ST_NATIONAL SelectorType = iota
	ST_REGIONAL
	ST_SEARCH
)

type Selector interface {
	PushCurrentArticle() tea.Cmd
	SelectorType() SelectorType
	SetActive(bool)
	IsActive() bool
	SetFocused(bool)
	IsFocused() bool
	SetDims(int, int)
	Init() tea.Cmd
	Update(tea.Msg) (Selector, tea.Cmd)
	View() string
}

type BaseSelector struct {
	selectorType  SelectorType
	shared        *SharedState
	isActive      bool
	isFocused     bool
	width         int
	height        int
	articles      []tagesschau.Article
	list          list.Model
	selectedIndex int
}

func NewSelector(selectorType SelectorType, shared *SharedState, isActive bool) BaseSelector {
	listKeymap := ListKeymap(shared.keys)
	return BaseSelector{
		shared:        shared,
		selectorType:  selectorType,
		isActive:      isActive,
		isFocused:     isActive,
		list:          initList(shared.style, listKeymap),
		selectedIndex: 0,
	}
}
func initList(s config.Style, km list.KeyMap) list.Model {
	lst := list.New([]list.Item{}, NewNewsDelegate(s), 0, 0)
	lst.SetFilteringEnabled(false)
	lst.SetShowTitle(false)
	lst.SetShowStatusBar(false)
	lst.SetShowHelp(false)
	lst.KeyMap = km
	return lst
}

func (s *BaseSelector) rebuildList() {
	var items []list.Item
	for _, article := range s.articles {
		items = append(items, article)
	}
	s.list.SetItems(items)
	s.list.Select(0)
	s.selectedIndex = 0
}

func (s *BaseSelector) getSelectedArticle() tagesschau.Article {
	return s.articles[s.selectedIndex]
}

func (s BaseSelector) PushCurrentArticle() tea.Cmd {
	return func() tea.Msg {
		return ChangedActiveArticle(s.getSelectedArticle())
	}
}

func (s BaseSelector) SelectorType() SelectorType {
	return s.selectorType
}

func (s *BaseSelector) SetActive(isActive bool) {
	s.isActive = isActive
}

func (s *BaseSelector) IsActive() bool {
	return s.isActive
}

func (s *BaseSelector) SetFocused(isFocused bool) {
	s.isFocused = isFocused
}

func (s *BaseSelector) IsFocused() bool {
	return s.isFocused
}

func (s *BaseSelector) SetDims(w, h int) {
	s.width = w
	s.height = h
}

func (s BaseSelector) Init() tea.Cmd {
	return nil
}

func (s BaseSelector) Update(msg tea.Msg) (BaseSelector, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.shared.mode == INSERT_MODE {
			break
		}

		switch {
		case key.Matches(msg, s.shared.keymap.right):
			if s.isActive {
				s.isFocused = false
			}
		case key.Matches(msg, s.shared.keymap.left):
			if s.isActive {
				s.isFocused = true
			}
		}
	}

	if s.list.Index() != s.selectedIndex {
		s.selectedIndex = s.list.Index()
		cmds = append(cmds, s.PushCurrentArticle())
	}

	return s, tea.Batch(cmds...)
}

func (s BaseSelector) View() string {
	s.list.SetSize(s.width, s.height)

	return s.list.View()
}
