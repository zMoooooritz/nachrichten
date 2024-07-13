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
type SwitchDirection int

const (
	HS_HIDDEN HelpState = iota
	HS_NORMAL
	HS_ALL

	SD_NEXT SwitchDirection = iota
	SD_PREV
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
	opener    util.Opener
	keymap    KeyMap
	style     config.Style
	ready     bool
	help      help.Model
	helpState HelpState
	selector  Selector
	viewers   []ViewerImplementation
	spinner   spinner.Model
	config    config.Configuration
	width     int
	height    int
}

func InitialModel(c config.Configuration) Model {
	style := config.NewsStyle(c.Theme)

	helpState := HS_NORMAL
	if c.Settings.HideHelpOnStartup {
		helpState = HS_HIDDEN
	}

	viewers := []ViewerImplementation{}
	viewers = append(viewers, NewReader(NewViewer(VT_TEXT, style, viewportKeymap(c.Keys), true)))
	viewers = append(viewers, NewImageViewer(NewViewer(VT_IMAGE, style, viewportKeymap(c.Keys), false)))

	m := Model{
		opener:    util.NewOpener(c.Applications),
		keymap:    GetKeyMap(c.Keys),
		style:     style,
		ready:     false,
		help:      NewHelper(style),
		helpState: helpState,
		selector:  NewSelector(style, listKeymap(c.Keys)),
		viewers:   viewers,
		spinner:   NewDotSpinner(),
		config:    c,
		width:     0,
		height:    0,
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
				m.switchViewer(SD_NEXT)
			}
		case key.Matches(msg, m.keymap.prev):
			if m.selector.IsFocused() {
				m.selector.PrevList()
				m.updateDisplayedArticle()
			} else {
				m.switchViewer(SD_PREV)
			}
		case key.Matches(msg, m.keymap.full):
			for idx1, viewer := range m.viewers {
				if viewer.IsActive() {
					viewer.SetFocused(true)
					m.selector.SetFocused(false)
					currentState := viewer.IsFullScreen()
					for idx2, viewer2 := range m.viewers {
						if idx1 == idx2 {
							viewer2.SetFullScreen(!currentState)
						} else {
							viewer2.SetFullScreen(currentState)
						}
					}
				}
			}

			m.updateSizes(m.width, m.height)
			m.updateDisplayedArticle()
		case key.Matches(msg, m.keymap.start):
			for _, viewer := range m.viewers {
				if viewer.IsFocused() {
					viewer.GotoTop()
				}
			}
		case key.Matches(msg, m.keymap.end):
			for _, viewer := range m.viewers {
				if viewer.IsFocused() {
					viewer.GotoBottom()
				}
			}
		case key.Matches(msg, m.keymap.image):
			m.showImageViewer()
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

	updatedViewer := false
	for i, viewer := range m.viewers {
		if viewer.IsFocused() || viewer.IsFullScreen() {
			updateViewer, cmd := viewer.Update(msg)
			m.viewers[i] = updateViewer
			cmds = append(cmds, cmd)
			updatedViewer = true
		}
	}
	if !updatedViewer {
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
		for _, viewer := range m.viewers {
			if viewer.IsActive() {
				viewer.SetFocused(true)
			}
		}
		m.selector.SetFocused(false)
	} else {
		for _, viewer := range m.viewers {
			viewer.SetFocused(false)
		}
		m.selector.SetFocused(true)
	}
	m.updateDisplayedArticle()
}

func (m *Model) switchViewer(switchDirection SwitchDirection) {
	currentIndex := -1
	for index, viewer := range m.viewers {
		if viewer.IsActive() {
			currentIndex = index
		}
	}
	newIndex := 0
	if currentIndex != -1 {
		if switchDirection == SD_NEXT {
			newIndex = (currentIndex + 1) % len(m.viewers)
		} else {
			newIndex = (currentIndex + len(m.viewers) - 1) % len(m.viewers)
		}

		m.viewers[currentIndex].SetActive(false)
		m.viewers[currentIndex].SetFocused(false)
		m.viewers[newIndex].SetActive(true)
		m.viewers[newIndex].SetFocused(true)
	}
	m.updateDisplayedArticle()
}

func (m *Model) showImageViewer() {
	for _, viewer := range m.viewers {
		if viewer.ViewerType() == VT_IMAGE {
			viewer.SetActive(true)
			viewer.SetFocused(true)
		} else {
			viewer.SetActive(false)
			viewer.SetFocused(false)
		}
	}
	m.selector.SetFocused(false)
}

func (m *Model) updateDisplayedArticle() {
	if !m.ready {
		return
	}

	article := m.getSelectedArticle()

	for _, viewer := range m.viewers {
		if viewer.IsActive() {
			if viewer.ViewerType() == VT_IMAGE {
				image := article.Thumbnail
				if image == nil {
					var err error
					image, err = http.LoadImage(article.ImageData.ImageVariants.RectSmall)
					if err != nil {
						return
					}
					article.Thumbnail = image
				}
			}
			viewer.SetArticle(*article)
			viewer.SetHeaderData(article.Topline, article.Date)
		}
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

	isViewerFullscreen := false
	for _, viewer := range m.viewers {
		if viewer.IsFullScreen() {
			m.selector.SetVisible(false)
			viewer.SetDims(m.width, m.height-m.helperHeight())
			isViewerFullscreen = true
		}
	}
	if !isViewerFullscreen {
		m.selector.SetVisible(true)
		for _, viewer := range m.viewers {
			viewer.SetDims(m.width-m.width/3-6, m.height-m.helperHeight())
		}
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
	viewerRepr := ""
	for _, viewer := range m.viewers {
		if viewer.IsActive() {
			viewerRepr = viewer.View()
		}
	}

	help := ""
	if m.helpState == HS_NORMAL || m.helpState == HS_ALL {
		help = "\n" + lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Center).Render(m.help.View(m.keymap))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, selector, viewerRepr) + help
}
