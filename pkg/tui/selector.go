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
	fullWidth          int
	headerWidth        int
	height             int
}

func NewSelector(s config.Style) Selector {
	return Selector{
		style:              s,
		lists:              InitLists(s, 2),
		listsActiveIndeces: []int{},
		activeListIndex:    0,
		isFocused:          true,
	}
}

func InitLists(s config.Style, count int) []list.Model {
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
		s.lists[i].SetSize(s.fullWidth, s.height)
		s.lists[i].Title = lipgloss.PlaceHorizontal(s.headerWidth, lipgloss.Center, headerText)
		s.lists[i].Styles.Title = s.style.TitleActiveStyle
	}
}

func (s *Selector) NextList() {
	s.activeListIndex = (s.activeListIndex + 1) % len(s.lists)
}

func (s *Selector) PrevList() {
	s.activeListIndex = (len(s.lists) + s.activeListIndex - 1) % len(s.lists)
}

func (s *Selector) SetFocused(isFocused bool) {
	s.isFocused = isFocused
}

func (s *Selector) SetDims(w, h int) {
	s.fullWidth = w
	s.headerWidth = w - 4
	s.height = h
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
	listHeader := s.listView([]string{nationalHeaderText, regionalHeaderText}, s.activeListIndex)
	listStyle := s.style.ListActiveStyle
	if !s.isFocused {
		listStyle = s.style.ListInactiveStyle
	}
	return listStyle.Render(lipgloss.JoinVertical(lipgloss.Left, listHeader, s.lists[s.activeListIndex].View()))
}

func (s Selector) listView(names []string, activeIndex int) string {
	cellWidth := s.headerWidth / len(names)
	var widths []int
	for i := 0; i < len(names)-1; i++ {
		widths = append(widths, cellWidth)
	}
	widths = append(widths, s.headerWidth-(len(names)-1)*cellWidth)
	result := ""
	for i, n := range names {
		style := s.style.TitleInactiveStyle
		if i == activeIndex {
			style = s.style.TitleActiveStyle
		}
		result += style.Render(lipgloss.PlaceHorizontal(widths[i], lipgloss.Center, n))
	}
	return lipgloss.NewStyle().PaddingLeft(2).Render(result)
}
