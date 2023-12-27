package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/nachrichten/pkg/config"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
)

const (
	headerText         string = "Nachrichten"
	regionalHeaderText string = "Regional"
	nationalHeaderText string = "National"
)

type Selector struct {
	style              config.Style
	news               tagesschau.News
	lists              []list.Model
	listsActiveIndeces []int
	activeListIndex    int
	isFocused          bool
	isVisible          bool
	fullWidth          int
	headerWidth        int
	fullHeight         int
	listHeight         int
}

func NewSelector(s config.Style) Selector {
	return Selector{
		style:              s,
		lists:              InitLists(s, 2),
		listsActiveIndeces: []int{},
		activeListIndex:    0,
		isFocused:          true,
		isVisible:          true,
	}
}

func InitLists(s config.Style, count int) []list.Model {
	var lists []list.Model
	for i := 0; i < count; i++ {
		newList := list.New([]list.Item{}, NewNewsDelegate(s), 0, 0)
		newList.SetFilteringEnabled(false)
		newList.SetShowTitle(false)
		newList.SetShowStatusBar(false)
		newList.SetShowHelp(false)
		lists = append(lists, newList)
	}
	return lists
}

func (s *Selector) FillLists(news tagesschau.News) {
	s.news = news
	for i, n := range [][]tagesschau.NewsEntry{news.NationalNews, news.RegionalNews} {
		var items []list.Item
		for _, ne := range n {
			items = append(items, ne)
		}

		s.lists[i].SetItems(items)
		s.listsActiveIndeces = append(s.listsActiveIndeces, 0)
	}
}

func (s *Selector) ResizeLists() {
	for i := range s.lists {
		s.lists[i].SetSize(s.fullWidth, s.listHeight)
	}
}

func (s *Selector) NextList() {
	s.activeListIndex = (s.activeListIndex + 1) % len(s.lists)
}

func (s *Selector) PrevList() {
	s.activeListIndex = (len(s.lists) + s.activeListIndex - 1) % len(s.lists)
}

func (s *Selector) IsFocused() bool {
	return s.isFocused
}

func (s *Selector) SetFocused(isFocused bool) {
	s.isFocused = isFocused
}

func (s *Selector) SetVisible(isVisible bool) {
	s.isVisible = isVisible
}

func (s *Selector) SetDims(w, h int) {
	s.fullWidth = w
	s.headerWidth = w - 2
	s.fullHeight = h
	s.listHeight = h - 4
}

func (s *Selector) GetSelectedArticle() tagesschau.NewsEntry {
	var article tagesschau.NewsEntry
	if s.activeListIndex == 0 {
		article = s.news.NationalNews[s.listsActiveIndeces[s.activeListIndex]]
	} else {
		article = s.news.RegionalNews[s.listsActiveIndeces[s.activeListIndex]]
	}
	return article
}

func (s *Selector) HasSelectionChanged() bool {
	if s.listsActiveIndeces[s.activeListIndex] != s.lists[s.activeListIndex].Index() {
		s.listsActiveIndeces[s.activeListIndex] = s.lists[s.activeListIndex].Index()
		return true
	}
	return false
}

func (s Selector) Init() tea.Cmd {
	return nil
}

func (s Selector) Update(msg tea.Msg) (Selector, tea.Cmd) {
	var cmd tea.Cmd
	s.lists[s.activeListIndex], cmd = s.lists[s.activeListIndex].Update(msg)
	return s, tea.Batch(cmd)
}

func (s Selector) View() string {
	if !s.isVisible {
		return ""
	}

	listSelect := s.listSelectView([]string{nationalHeaderText, regionalHeaderText}, s.activeListIndex)
	listHead := s.listHeadView()

	style := s.style.ListInactiveStyle
	if s.isFocused {
		style = s.style.ListActiveStyle
	}

	return style.Render(lipgloss.JoinVertical(lipgloss.Left, listHead, listSelect, s.lists[s.activeListIndex].View()))
}

func (s Selector) listHeadView() string {
	headerStyle := s.style.ListHeaderInactiveStyle
	if s.isFocused {
		headerStyle = s.style.ListHeaderActiveStyle
	}

	centeredHeader := lipgloss.PlaceHorizontal(s.headerWidth, lipgloss.Center, headerText)
	return headerStyle.Render(centeredHeader)
}

func (s Selector) listSelectView(names []string, activeIndex int) string {
	cellWidth := (s.headerWidth - len(names)) / len(names)
	var widths []int
	for i := 0; i < len(names)-1; i++ {
		widths = append(widths, cellWidth)
	}
	widths = append(widths, s.headerWidth-(len(names)-1)*cellWidth-len(names))
	result := ""
	for i, n := range names {
		border := s.style.InactiveTabBorder
		style := s.style.InactiveStyle
		if i == activeIndex {
			border = s.style.ActiveTabBorder
			style = s.style.TextHighlightStyle
		}
		centeredText := lipgloss.PlaceHorizontal(widths[i], lipgloss.Center, n)
		result = lipgloss.JoinHorizontal(lipgloss.Center, result, style.Copy().MarginBottom(1).BorderStyle(border).Render(centeredText))
	}
	return result
}
