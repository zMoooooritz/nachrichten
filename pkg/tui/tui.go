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

type HelpState int

const (
	HS_HIDDEN HelpState = iota
	HS_NORMAL
	HS_ALL
)

var (
	screenCentered = func(w, h int) lipgloss.Style {
		return lipgloss.NewStyle().
			Width(w).
			Align(lipgloss.Center).
			Height(h).
			AlignVertical(lipgloss.Center)
	}

	news      tagesschau.News
	imageSpec = tagesschau.ImageSpec{Size: tagesschau.SMALL, Ratio: tagesschau.RECT}
)

type Model struct {
	opener      util.Opener
	keymap      KeyMap
	style       config.Style
	ready       bool
	help        help.Model
	helpState   HelpState
	selector    Selector
	reader      Reader
	imageViewer ImageViewer
	spinner     spinner.Model
	config      config.Configuration
	width       int
	height      int
}

func InitialModel(c config.Configuration) Model {
	style := config.NewsStyle(c.Theme)

	helpState := HS_NORMAL
	if c.Settings.HideHelpOnStartup {
		helpState = HS_HIDDEN
	}

	m := Model{
		opener:      util.NewOpener(c.Applications),
		keymap:      GetKeyMap(c.Keys),
		style:       style,
		ready:       false,
		help:        NewHelper(style),
		helpState:   helpState,
		selector:    NewSelector(style, listKeymap(c.Keys)),
		reader:      NewReader(style, viewportKeymap(c.Keys), true),
		imageViewer: NewImageViewer(style, viewportKeymap(c.Keys), false),
		spinner:     NewDotSpinner(),
		config:      c,
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
		news = tagesschau.News(msg)
		if m.config.Settings.PreloadThumbnails {
			go news.EnrichWithThumbnails(imageSpec)
		}
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
			m.helpState = (m.helpState + 1) % 3
			if m.helpState == HS_NORMAL {
				m.help.ShowAll = false
			}
			if m.helpState == HS_ALL {
				m.help.ShowAll = true
			}
			m.updateSizes(m.width, m.height)
		case key.Matches(msg, m.keymap.right):
			m.setFocus(true)
		case key.Matches(msg, m.keymap.left):
			m.setFocus(false)
		case key.Matches(msg, m.keymap.next):
			if m.selector.IsFocused() {
				m.selector.NextList()
				m.updateDisplayedArticle()
			} else {
				m.toggelViewer(true)
			}
		case key.Matches(msg, m.keymap.prev):
			if m.selector.IsFocused() {
				m.selector.PrevList()
				m.updateDisplayedArticle()
			} else {
				m.toggelViewer(true)
			}
		case key.Matches(msg, m.keymap.full):
			if m.reader.IsActive() {
				m.reader.SetFocused(true)
				m.selector.SetFocused(false)
				m.reader.SetFullScreen(!m.reader.IsFullScreen())
				m.imageViewer.SetFullScreen(m.reader.IsFullScreen())
			}
			if m.imageViewer.IsActive() {
				m.imageViewer.SetFocused(true)
				m.selector.SetFocused(false)
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
			m.toggelViewer(!m.selector.IsFocused())
		case key.Matches(msg, m.keymap.open):
			article := m.getSelectedArticle()
			m.opener.OpenUrl(util.TypeHTML, article.URL)
		case key.Matches(msg, m.keymap.video):
			article := m.getSelectedArticle()
			m.opener.OpenUrl(util.TypeVideo, article.Video.VideoVariants.Big)
		case key.Matches(msg, m.keymap.shortNews):
			url, err := tagesschau.GetShortNewsURL()
			if err == nil {
				m.opener.OpenUrl(util.TypeVideo, url)
			}
		}
	case tea.WindowSizeMsg:
		m.updateSizes(msg.Width, msg.Height)
		m.updateDisplayedArticle()
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

func (m *Model) toggelViewer(setFocus bool) {
	m.reader.SetActive(!m.reader.IsActive())
	m.imageViewer.SetActive(!m.imageViewer.IsActive())
	if setFocus {
		m.reader.SetFocused(!m.reader.IsFocused())
		m.imageViewer.SetFocused(!m.imageViewer.IsFocused())
	}
	m.updateDisplayedArticle()
}

func (m *Model) updateDisplayedArticle() {
	if !m.ready {
		return
	}

	article := m.getSelectedArticle()
	if m.reader.isActive {
		m.reader.SetArticle(*article)
		m.reader.SetHeaderData(article.Topline, article.Date)
	}

	if m.imageViewer.IsActive() {
		image := article.Thumbnail
		if image == nil {
			var err error
			image, err = http.LoadImage(article.ImageData.ImageVariants.RectSmall)
			if err != nil {
				return
			}
			article.Thumbnail = image
		}
		m.imageViewer.SetArticle(*article)
		m.imageViewer.SetHeaderData(article.Topline, article.Date)
	}
}

func (m *Model) getSelectedArticle() *tagesschau.Article {
	newsIndex, articleIndex := m.selector.GetSelectedIndex()
	if newsIndex == 0 {
		return &news.NationalNews[articleIndex]
	}
	return &news.RegionalNews[articleIndex]
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
	if m.helpState == HS_NORMAL || m.helpState == HS_ALL {
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
	if m.helpState == HS_NORMAL || m.helpState == HS_ALL {
		help = "\n" + lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Center).Render(m.help.View(m.keymap))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, selector, viewer) + help
}
