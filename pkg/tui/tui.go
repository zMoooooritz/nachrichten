package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/nachrichten/pkg/config"
	"github.com/zMoooooritz/nachrichten/pkg/http"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
	"github.com/zMoooooritz/nachrichten/pkg/util"
)

const (
	germanDateFormat string = "15:04 02.01.06"
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
	opener      util.Opener
	keymap      KeyMap
	style       config.Style
	ready       bool
	help        help.Model
	helpMode    int
	selector    Selector
	reader      Reader
	imageViewer ImageViewer
	spinner     spinner.Model
	width       int
	height      int
}

func InitialModel(c config.Configuration) Model {
	theme := config.Theme{}
	var style config.Style
	if c.Theme != theme {
		style = config.NewsStyle(c.Theme)
	} else {
		style = config.NewsStyle(config.GruvboxTheme())
	}

	helpMode := 1
	if c.SettingsConfig.HideHelpOnStartup {
		helpMode = 0
	}

	m := Model{
		opener:      util.NewOpener(c),
		keymap:      GetKeyMap(),
		style:       style,
		ready:       false,
		help:        NewHelper(style),
		helpMode:    helpMode,
		selector:    NewSelector(style),
		reader:      NewReader(style),
		imageViewer: NewImageViewer(style),
		spinner:     NewDotSpinner(),
		width:       0,
		height:      0,
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
		news := tagesschau.News(msg)
		m.selector.FillLists(news)
		m.selector.ResizeLists()
		m.ready = true
		m.updateDisplayedArticle()
	case tea.KeyMsg:
		if !m.ready {
			return m, tea.Batch(cmds...)
		}

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
			m.setFocus(true)
		case key.Matches(msg, m.keymap.left):
			m.setFocus(false)
		case key.Matches(msg, m.keymap.next):
			m.selector.NextList()
			m.setFocus(false)
		case key.Matches(msg, m.keymap.prev):
			m.selector.PrevList()
			m.setFocus(false)
		case key.Matches(msg, m.keymap.full):
			if m.reader.IsActive() {
				m.reader.SetFullScreen(!m.reader.IsFullScreen())
				m.imageViewer.SetFullScreen(m.reader.IsFullScreen())
			}
			if m.imageViewer.IsActive() {
				m.imageViewer.SetFullScreen(!m.imageViewer.IsFullScreen())
				m.reader.SetFullScreen(m.imageViewer.IsFullScreen())
			}
			m.updateSizes(m.width, m.height)
			m.updateDisplayedArticle()
		case key.Matches(msg, m.keymap.start):
			if m.reader.IsFocused() {
				m.reader.GotoTop()
			}
		case key.Matches(msg, m.keymap.end):
			if m.reader.IsFocused() {
				m.reader.GotoBottom()
			}
		case key.Matches(msg, m.keymap.image):
			if m.selector.IsFocused() {
				m.reader.SetActive(!m.reader.IsActive())
				m.imageViewer.SetActive(!m.imageViewer.IsActive())
				m.updateDisplayedArticle()
				// m.updateSizes(m.width, m.height)
			}
		case key.Matches(msg, m.keymap.open):
			article := m.selector.GetSelectedArticle()
			m.opener.OpenUrl(config.TypeHTML, article.URL)
		case key.Matches(msg, m.keymap.video):
			article := m.selector.GetSelectedArticle()
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

	if m.reader.IsFocused() || m.reader.IsFullScreen() {
		m.reader, cmd = m.reader.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.imageViewer.IsFocused() || m.imageViewer.IsFullScreen() {
		m.imageViewer, cmd = m.imageViewer.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.selector, cmd = m.selector.Update(msg)
		cmds = append(cmds, cmd)
		if m.selector.HasSelectionChanged() {
			m.updateDisplayedArticle()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) setFocus(onViewer bool) {
	if onViewer {
		if m.reader.IsActive() {
			m.reader.SetFocused(true)
		}
		if m.imageViewer.IsActive() {
			m.imageViewer.SetFocused(true)
		}
		m.selector.SetFocused(false)
	} else {
		m.reader.SetFocused(false)
		m.imageViewer.SetFocused(false)
		m.selector.SetFocused(true)
	}
	m.updateDisplayedArticle()
}

func (m *Model) updateDisplayedArticle() {
	article := m.selector.GetSelectedArticle()
	if m.reader.isActive {
		text := tagesschau.ContentToParagraphs(article.Content)
		m.reader.SetContent(text)
		m.reader.SetHeaderContent(article.Topline, article.Date)
	}
	if m.imageViewer.IsActive() {
		image, err := http.LoadImage(article.Image.ImageURLs.RectSmall)
		if err != nil {
			return
		}
		m.imageViewer.SetImage(image)
		m.imageViewer.SetHeaderContent(article.Topline, article.Date)
	}
}

func (m *Model) updateSizes(width, height int) {
	m.width = width
	m.height = height

	m.selector.SetDims(m.width/3, m.height-m.helperHeight()-5)
	m.selector.ResizeLists()

	if m.reader.IsFullScreen() {
		m.selector.SetVisible(false)
		m.reader.SetDims(m.width, m.height-m.helperHeight())
	}
	if m.imageViewer.IsFullScreen() {
		m.selector.SetVisible(false)
		m.imageViewer.SetDims(m.width, m.height-m.helperHeight())
	}
	if !m.reader.IsFullScreen() && !m.imageViewer.IsFullScreen() {
		m.selector.SetVisible(true)
		m.reader.SetDims(m.width-m.width/3-6, m.height-m.helperHeight())
		m.imageViewer.SetDims(m.width-m.width/3-6, m.height-m.helperHeight())
	}

	m.help.Width = m.width
}

func (m Model) helperHeight() int {
	if m.helpMode > 0 {
		return 2
	}
	return 0
}

func (m Model) View() string {
	if !m.ready {
		content := fmt.Sprintf("%s Lade Nachrichten... press q to quit", m.spinner.View())
		return screenCentered(m.width, m.height).Render(content)
	}

	selector := m.selector.View()
	viewer := ""
	if m.reader.IsActive() {
		viewer = m.reader.View()
	}
	if m.imageViewer.IsActive() {
		viewer = m.imageViewer.View()
	}

	help := ""
	if m.helpMode > 0 {
		help = "\n" + lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Center).Render(m.help.View(m.keymap))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, selector, viewer) + help
}
