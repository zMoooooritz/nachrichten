package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/Nachrichten/pkg/config"
	"github.com/zMoooooritz/Nachrichten/pkg/tagesschau"
	"github.com/zMoooooritz/Nachrichten/pkg/util"
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
	configuration      config.Configuration
	news               tagesschau.News
	keymap             KeyMap
	style              config.Style
	ready              bool
	help               help.Model
	helpMode           int
	lists              []list.Model
	listsActiveIndeces []int
	activeListIndex    int
	reader             viewport.Model
	spinner            spinner.Model
	focus              int
	readerFocused      bool
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

	m := Model{
		configuration:      c,
		keymap:             GetKeyMap(),
		style:              style,
		ready:              false,
		help:               NewHelper(style),
		helpMode:           1,
		reader:             viewport.New(0, 0),
		spinner:            NewDotSpinner(),
		focus:              0,
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
		for i := range m.lists {
			width, _ := m.listInnerDims()
			m.lists[i].Title = lipgloss.PlaceHorizontal(width, lipgloss.Center, headerText)
		}
		m.ready = true
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
			m.readerFocused = true
		case key.Matches(msg, m.keymap.left):
			m.readerFocused = false
		case key.Matches(msg, m.keymap.next):
			m.readerFocused = false
			m.activeListIndex = (m.activeListIndex + 1) % len(m.lists)
		case key.Matches(msg, m.keymap.start):
			if m.readerFocused {
				m.reader.GotoTop()
			}
		case key.Matches(msg, m.keymap.end):
			if m.readerFocused {
				m.reader.GotoBottom()
			}
		case key.Matches(msg, m.keymap.prev):
			m.readerFocused = false
			m.activeListIndex = (len(m.lists) + m.activeListIndex - 1) % len(m.lists)
		case key.Matches(msg, m.keymap.open):
			article := m.SelectedArticle()
			_ = util.OpenUrl(config.TypeHTML, m.configuration, article.URL)
		case key.Matches(msg, m.keymap.video):
			article := m.SelectedArticle()
			_ = util.OpenUrl(config.TypeVideo, m.configuration, article.Video.VideoURLs.Big)
		case key.Matches(msg, m.keymap.shortNews):
			url, _ := tagesschau.GetShortNewsURL()
			_ = util.OpenUrl(config.TypeVideo, m.configuration, url)
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

	if m.readerFocused {
		m.reader, cmd = m.reader.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.lists[m.activeListIndex], cmd = m.lists[m.activeListIndex].Update(msg)
		cmds = append(cmds, cmd)
		m.listsActiveIndeces[m.activeListIndex] = m.lists[m.activeListIndex].Index()
		m.reader.SetContent(util.ContentToText(m.SelectedArticle().Content, m.reader.Width))
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateSizes(width, height int) {
	m.width = width
	m.height = height

	m.reader.YPosition = m.readerHeaderHeight()

	for i := range m.lists {
		m.lists[i].SetSize(m.listOuterDims())
		w, _ := m.listInnerDims()
		m.lists[i].Title = lipgloss.PlaceHorizontal(w, lipgloss.Center, headerText)
	}

	m.reader.Width, m.reader.Height = m.readerDims()
	m.help.Width = m.width
}

func (m Model) listOuterDims() (int, int) {
	return m.width / 3, m.height - m.helperHeight() - 5
}

func (m Model) listInnerDims() (int, int) {
	w, h := m.listOuterDims()
	return w - 6, h
}

func (m Model) listSelectorDims() (int, int) {
	w, h := m.listOuterDims()
	return w - 4, h
}

func (m Model) readerDims() (int, int) {
	lw, _ := m.listOuterDims()
	return m.width - lw - 7, m.height - m.readerHeaderHeight() - m.readerFooterHeight() - m.helperHeight()
}

func (m Model) readerHeaderHeight() int {
	return lipgloss.Height(m.headerView("", ""))
}

func (m Model) readerFooterHeight() int {
	return lipgloss.Height(m.footerView())
}

func (m Model) helperHeight() int {
	if m.helpMode > 0 {
		return 2
	}
	return 0
}

func (m Model) SelectedArticle() tagesschau.NewsEntry {
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

	listHeader := m.listSelectorView([]string{nationalHeaderText, regionalHeaderText}, m.activeListIndex)
	listStyle := m.style.ListActiveStyle
	if m.readerFocused {
		listStyle = m.style.ListInactiveStyle
	}
	list := listStyle.Render(lipgloss.JoinVertical(lipgloss.Left, listHeader, m.lists[m.activeListIndex].View()))
	article := m.SelectedArticle()
	reader := fmt.Sprintf("%s\n%s\n%s", m.headerView(article.Topline, article.Date.Format(germanDateFormat)), m.reader.View(), m.footerView())

	help := ""
	if m.helpMode > 0 {
		help = "\n" + lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Center).Render(m.help.View(m.keymap))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, list, reader) + help
}

func (m Model) listSelectorView(names []string, activeIndex int) string {
	width, _ := m.listSelectorDims()
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

func (m Model) headerView(name string, date string) string {
	titleStyle := m.style.ReaderTitleInactiveStyle
	lineStyle := m.style.InactiveStyle
	dateStyle := m.style.ReaderInfoInactiveStyle
	if m.readerFocused {
		titleStyle = m.style.ReaderTitleActiveStyle
		lineStyle = m.style.ActiveStyle
		dateStyle = m.style.ReaderInfoActiveStyle
	}

	title := titleStyle.Render(name)
	date = dateStyle.Render(date)
	line := lineStyle.Render(strings.Repeat("─", util.Max(0, m.reader.Width-lipgloss.Width(title)-lipgloss.Width(date))))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line, date)
}

func (m Model) footerView() string {
	infoStyle := m.style.ReaderInfoInactiveStyle
	lineStyle := m.style.InactiveStyle
	if m.readerFocused {
		infoStyle = m.style.ReaderInfoActiveStyle
		lineStyle = m.style.ActiveStyle
	}

	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.reader.ScrollPercent()*100))
	line := lineStyle.Render(strings.Repeat("─", util.Max(0, m.reader.Width-lipgloss.Width(info))))

	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
