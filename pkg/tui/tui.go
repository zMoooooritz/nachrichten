package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/nachrichten/pkg/config"
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
	opener        util.Opener
	keymap        KeyMap
	style         config.Style
	ready         bool
	help          help.Model
	helpState     HelpState
	selector      Selector
	viewers       []Viewer
	spinner       spinner.Model
	config        config.Configuration
	activeArticle *tagesschau.Article
	imageCache    *ImageCache
	width         int
	height        int
}

func InitialModel(c config.Configuration) Model {
	style := config.NewsStyle(c.Theme)

	helpState := HS_NORMAL
	if c.Settings.HideHelpOnStartup {
		helpState = HS_HIDDEN
	}

	ic := NewImageCache()

	viewers := []Viewer{}
	viewers = append(viewers, NewReader(NewViewer(VT_TEXT, style, viewportKeymap(c.Keys), true)))
	viewers = append(viewers, NewImageViewer(NewViewer(VT_IMAGE, style, viewportKeymap(c.Keys), false), ic))
	viewers = append(viewers, NewDetails(NewViewer(VT_DETAILS, style, viewportKeymap(c.Keys), false)))

	m := Model{
		opener:        util.NewOpener(c.Applications),
		keymap:        GetKeyMap(c.Keys),
		style:         style,
		ready:         false,
		help:          NewHelper(style),
		helpState:     helpState,
		selector:      NewSelector(style, listKeymap(c.Keys)),
		viewers:       viewers,
		spinner:       NewDotSpinner(),
		config:        c,
		activeArticle: nil,
		imageCache:    ic,
		width:         0,
		height:        0,
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
			go func() {
				for _, a := range news.NationalNews {
					_ = m.imageCache.LoadImage(a.ID, a.ImageData.ImageVariants.RectSmall)
				}
				for _, a := range news.RegionalNews {
					_ = m.imageCache.LoadImage(a.ID, a.ImageData.ImageVariants.RectSmall)
				}
			}()
		}
		m.selector.FillLists(news)
		m.selector.ResizeLists()
		m.ready = true
		m.activeArticle = &news.NationalNews[0]
		m.updateActiveViewer()
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
				m.updateActiveArticle()
				m.updateActiveViewer()
			}
		case key.Matches(msg, m.keymap.prev):
			if m.selector.IsFocused() {
				m.selector.PrevList()
				m.updateActiveArticle()
				m.updateActiveViewer()
			}
		case key.Matches(msg, m.keymap.full):
			activeViewer := m.activeViewer()
			activeViewer.SetFocused(true)
			m.selector.SetFocused(false)
			currentState := activeViewer.IsFullScreen()
			for _, viewer := range m.viewers {
				if activeViewer.ViewerType() == activeViewer.ViewerType() {
					viewer.SetFullScreen(!currentState)
				} else {
					viewer.SetFullScreen(currentState)
				}
			}

			m.updateSizes(m.width, m.height)
			m.updateActiveViewer()
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
		case key.Matches(msg, m.keymap.article):
			m.showViewer(VT_TEXT)
		case key.Matches(msg, m.keymap.image):
			m.showViewer(VT_IMAGE)
		case key.Matches(msg, m.keymap.details):
			m.showViewer(VT_DETAILS)
		case key.Matches(msg, m.keymap.open):
			m.opener.OpenUrl(util.TypeHTML, m.activeArticle.URL)
		case key.Matches(msg, m.keymap.video):
			m.opener.OpenUrl(util.TypeVideo, m.activeArticle.Video.VideoVariants.Big)
		case key.Matches(msg, m.keymap.shortNews):
			url, err := tagesschau.GetShortNewsURL()
			if err == nil {
				m.opener.OpenUrl(util.TypeVideo, url)
			}
		}
		keyStr := msg.String()
		if keyStr >= "0" && keyStr <= "9" {
			keyInt, _ := strconv.Atoi(keyStr)
			m.handleNumberInput(keyInt)
		}
	case tea.WindowSizeMsg:
		m.updateSizes(msg.Width, msg.Height)
		m.updateActiveViewer()
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
			m.updateActiveArticle()
			m.updateActiveViewer()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) activeViewer() Viewer {
	for _, viewer := range m.viewers {
		if viewer.IsActive() {
			return viewer
		}
	}
	return m.viewers[0]
}

func (m *Model) handleNumberInput(number int) {
	if m.activeViewer().ViewerType() == VT_DETAILS {
		related := m.activeArticle.GetRelatedArticles()
		index := number - 1
		if 0 <= index && index < len(related) {
			m.activeArticle = tagesschau.LoadArticle(related[index].Details)
			m.showViewer(VT_TEXT)
			m.updateActiveViewer()
		}
	}
}

func (m *Model) setFocus(onViewer bool) {
	if onViewer {
		m.activeViewer().SetFocused(true)
		m.selector.SetFocused(false)
	} else {
		for _, viewer := range m.viewers {
			viewer.SetFocused(false)
		}
		m.selector.SetFocused(true)
	}
	m.updateActiveViewer()
}

func (m *Model) showViewer(vt ViewerType) {
	isViewerFocused := m.activeViewer().IsFocused()
	for _, viewer := range m.viewers {
		if viewer.ViewerType() == vt {
			viewer.SetActive(true)
			viewer.SetFocused(isViewerFocused)
		} else {
			viewer.SetActive(false)
			viewer.SetFocused(false)
		}
	}
	m.updateActiveViewer()
}

func (m *Model) updateActiveViewer() {
	if !m.ready {
		return
	}

	m.activeViewer().SetArticle(*m.activeArticle)
}

func (m *Model) updateActiveArticle() {
	if m.selector.IsFocused() {
		newsIndex, articleIndex := m.selector.GetSelectedIndex()
		if newsIndex == 0 {
			m.activeArticle = &news.NationalNews[articleIndex]
		} else {
			m.activeArticle = &news.RegionalNews[articleIndex]
		}
	}
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
	viewer := m.activeViewer().View()
	help := ""
	if m.helpState == HS_NORMAL || m.helpState == HS_ALL {
		help = "\n" + lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Center).Render(m.help.View(m.keymap))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, selector, viewer) + help
}
