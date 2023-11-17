package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/nachrichten/pkg/config"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
	"github.com/zMoooooritz/nachrichten/pkg/util"
)

const (
	headerText         string = "Nachrichten"
	regionalHeaderText string = "Regional"
	nationalHeaderText string = "National"
	germanDateFormat   string = "15:04 02.01.06"
)

var (
	screenCentered = func(w, h int) lipgloss.Style {
		return lipgloss.NewStyle().
			Width(w).
			Align(lipgloss.Center).
			Height(h).
			AlignVertical(lipgloss.Center)
	}
)

type Model struct {
	opener             util.Opener
	news               tagesschau.News
	keymap             KeyMap
	style              config.Style
	ready              bool
	help               help.Model
	helpMode           int
	lists              []list.Model
	listsActiveIndeces []int
	activeListIndex    int
	reader             Reader
	spinner            spinner.Model
	width              int
	height             int
}

func (m *Model) InitLists(news [][]tagesschau.NewsEntry) {
	for i, n := range news {
		var items []list.Item
		for _, ne := range n {
			items = append(items, ne)
		}

		m.lists[i].SetItems(items)
		m.listsActiveIndeces = append(m.listsActiveIndeces, 0)
	}
}

func InitialModel(c config.Configuration) Model {
	tc := config.ThemeConfig{}
	var style config.Style
	if c.ThemeConfig != tc {
		style = config.NewsStyle(c.ThemeConfig)
	} else {
		style = config.NewsStyle(config.DefaultThemeConfiguration())
	}

	helpMode := 1
	if c.SettingsConfig.HideHelpOnStartup {
		helpMode = 0
	}

	m := Model{
		opener:             util.NewOpener(c),
		keymap:             GetKeyMap(),
		style:              style,
		ready:              false,
		help:               NewHelper(style),
		helpMode:           helpMode,
		reader:             NewReader(style),
		spinner:            NewDotSpinner(),
		lists:              EmptyLists(style, 2),
		listsActiveIndeces: []int{},
		activeListIndex:    0,
		width:              0,
		height:             0,
	}
	return m
}

func GetNews() tea.Cmd {
	return func() tea.Msg {
		return tagesschau.LoadNews()
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(GetNews(), m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tagesschau.News:
		m.news = tagesschau.News(msg)
		m.InitLists([][]tagesschau.NewsEntry{m.news.NationalNews, m.news.RegionalNews})
		m.resizeLists()
		m.ready = true
		m.updateDisplayedArticle()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.help):
			m.helpMode = (m.helpMode + 1) % 3
			if m.helpMode == 1 {
				m.help.ShowAll = false
			}
			if m.helpMode == 2 {
				m.help.ShowAll = true
			}
			m.updateSizes(m.width, m.height)
		case key.Matches(msg, m.keymap.right):
			m.reader.SetFocused(true)
		case key.Matches(msg, m.keymap.left):
			m.reader.SetFocused(false)
		case key.Matches(msg, m.keymap.next):
			m.reader.SetFocused(false)
			m.activeListIndex = (m.activeListIndex + 1) % len(m.lists)
			m.updateDisplayedArticle()
		case key.Matches(msg, m.keymap.prev):
			m.reader.SetFocused(false)
			m.activeListIndex = (len(m.lists) + m.activeListIndex - 1) % len(m.lists)
			m.updateDisplayedArticle()
		case key.Matches(msg, m.keymap.start):
			if m.reader.IsFocused() {
				m.reader.GotoTop()
			}
		case key.Matches(msg, m.keymap.end):
			if m.reader.IsFocused() {
				m.reader.GotoBottom()
			}
		case key.Matches(msg, m.keymap.open):
			article := m.selectedArticle()
			m.opener.OpenUrl(config.TypeHTML, article.URL)
		case key.Matches(msg, m.keymap.video):
			article := m.selectedArticle()
			m.opener.OpenUrl(config.TypeVideo, article.Video.VideoURLs.Big)
		case key.Matches(msg, m.keymap.shortNews):
			url, _ := tagesschau.GetShortNewsURL()
			m.opener.OpenUrl(config.TypeVideo, url)
		}
	case tea.WindowSizeMsg:
		m.updateSizes(msg.Width, msg.Height)
	default:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if !m.ready {
		return m, tea.Batch(cmds...)
	}

	if m.reader.IsFocused() {
		m.reader, cmd = m.reader.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.lists[m.activeListIndex], cmd = m.lists[m.activeListIndex].Update(msg)
		cmds = append(cmds, cmd)
		if m.listsActiveIndeces[m.activeListIndex] != m.lists[m.activeListIndex].Index() {
			m.listsActiveIndeces[m.activeListIndex] = m.lists[m.activeListIndex].Index()
			m.updateDisplayedArticle()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateDisplayedArticle() {
	article := m.selectedArticle()
	text := tagesschau.ContentToParagraphs(article.Content)
	m.reader.SetContent(text)
	m.reader.SetHeaderContent(article.Topline, article.Date.Format(germanDateFormat))
}

func (m *Model) updateSizes(width, height int) {
	m.width = width
	m.height = height

	m.resizeLists()

	w, _ := m.listOuterDims()
	m.reader.SetDims(m.width-w-6, m.height-m.helperHeight())
	m.help.Width = m.width
}

func (m *Model) resizeLists() {
	w, _ := m.listInnerDims()
	for i := range m.lists {
		m.lists[i].SetSize(m.listOuterDims())
		m.lists[i].Title = lipgloss.PlaceHorizontal(w, lipgloss.Center, headerText)
		m.lists[i].Styles.Title = m.style.TitleActiveStyle
	}
}

func (m Model) listOuterDims() (int, int) {
	return m.width / 3, m.height - m.helperHeight() - 5
}

func (m Model) listInnerDims() (int, int) {
	w, h := m.listOuterDims()
	return w - 4, h
}

func (m Model) helperHeight() int {
	if m.helpMode > 0 {
		return 2
	}
	return 0
}

func (m Model) selectedArticle() tagesschau.NewsEntry {
	var article tagesschau.NewsEntry
	if m.activeListIndex == 0 {
		article = m.news.NationalNews[m.listsActiveIndeces[m.activeListIndex]]
	} else {
		article = m.news.RegionalNews[m.listsActiveIndeces[m.activeListIndex]]
	}
	return article
}

func (m Model) View() string {
	if !m.ready {
		content := fmt.Sprintf("%s Lade Nachrichten... press q to quit", m.spinner.View())
		return screenCentered(m.width, m.height).Render(content)
	}

	listHeader := m.listView([]string{nationalHeaderText, regionalHeaderText}, m.activeListIndex)
	listStyle := m.style.ListActiveStyle
	if m.reader.IsFocused() {
		listStyle = m.style.ListInactiveStyle
	}
	list := listStyle.Render(lipgloss.JoinVertical(lipgloss.Left, listHeader, m.lists[m.activeListIndex].View()))
	reader := m.reader.View()

	help := ""
	if m.helpMode > 0 {
		help = "\n" + lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Center).Render(m.help.View(m.keymap))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, list, reader) + help
}

func (m Model) listView(names []string, activeIndex int) string {
	width, _ := m.listInnerDims()
	cellWidth := width / len(names)
	var widths []int
	for i := 0; i < len(names)-1; i++ {
		widths = append(widths, cellWidth)
	}
	widths = append(widths, width-(len(names)-1)*cellWidth)
	result := ""
	for i, n := range names {
		style := m.style.TitleInactiveStyle
		if i == activeIndex {
			style = m.style.TitleActiveStyle
		}
		result += style.Render(lipgloss.PlaceHorizontal(widths[i], lipgloss.Center, n))
	}
	return lipgloss.NewStyle().PaddingLeft(2).Render(result)
}
