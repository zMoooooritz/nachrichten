package tui

import (
	"fmt"
	"strconv"

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

var (
	news tagesschau.News
)

type Model struct {
	opener        util.Opener
	keymap        KeyMap
	style         config.Style
	ready         bool
	loadingFailed bool
	helper        *Helper
	selector      *Selector
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

	initialHelpState := HS_NORMAL
	if c.Settings.HideHelpOnStartup {
		initialHelpState = HS_HIDDEN
	}

	ic := NewImageCache()

	viewers := []Viewer{
		NewReader(NewViewer(VT_TEXT, style, c.Keys, true)),
		NewImageViewer(NewViewer(VT_IMAGE, style, c.Keys, false), ic),
		NewDetails(NewViewer(VT_DETAILS, style, c.Keys, false)),
	}

	return Model{
		opener:     util.NewOpener(c.Applications),
		keymap:     GetKeyMap(c.Keys),
		style:      style,
		ready:      false,
		helper:     NewHelper(style, c.Keys, initialHelpState),
		selector:   NewSelector(style, c.Keys),
		viewers:    viewers,
		spinner:    NewDotSpinner(),
		config:     c,
		imageCache: ic,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(loadNews, m.spinner.Tick)
}

func loadNews() tea.Msg {
	news, err := tagesschau.LoadNews()
	if err == nil {
		return news
	}
	return LoadingNewsFailed{}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case LoadingNewsFailed:
		m.loadingFailed = true
	case tagesschau.News:
		news = tagesschau.News(msg)
		if m.config.Settings.PreloadThumbnails {
			go m.loadThumbnails()
		}
		m.ready = true
		m.activeArticle = &news.NationalNews[0]
		m.updateActiveViewer()
	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.quit) {
			return m, tea.Quit
		}

		if !m.ready {
			return m, nil
		}

		switch {
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
	}

	if !m.ready {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	m.helper, cmd = m.helper.Update(msg)
	cmds = append(cmds, cmd)

	m.selector, cmd = m.selector.Update(msg)
	cmds = append(cmds, cmd)
	if m.selector.HasSelectionChanged() {
		m.updateActiveArticle()
		m.updateActiveViewer()
	}

	for i, viewer := range m.viewers {
		updatedViewer, cmd := viewer.Update(msg)
		m.viewers[i] = updatedViewer
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.help):
			m.updateSizes(m.width, m.height)
		case key.Matches(msg, m.keymap.full):
			m.updateSizes(m.width, m.height)
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

func (m *Model) loadThumbnails() {
	imageSpec := tagesschau.ImageSpec{Size: tagesschau.SMALL, Ratio: tagesschau.RECT}
	for _, a := range news.NationalNews {
		_ = m.imageCache.LoadImage(a.ID, tagesschau.GetImageURL(a.ImageData.ImageVariants, imageSpec))
	}
	for _, a := range news.RegionalNews {
		_ = m.imageCache.LoadImage(a.ID, tagesschau.GetImageURL(a.ImageData.ImageVariants, imageSpec))
	}
}

func (m *Model) handleNumberInput(number int) {
	if m.activeViewer().ViewerType() == VT_DETAILS {
		related := m.activeArticle.GetRelatedArticles()
		index := number - 1
		if 0 <= index && index < len(related) {
			article, err := tagesschau.LoadArticle(related[index].Details)
			if err == nil {
				m.activeArticle = article
				m.showViewer(VT_TEXT)
				m.updateActiveViewer()
			}
		}
	}
}

func (m *Model) showViewer(vt ViewerType) {
	currViewer := m.activeViewer()
	nextViewer := m.activeViewer()
	for _, viewer := range m.viewers {
		if viewer.ViewerType() == vt {
			nextViewer = viewer
		}
	}
	if currViewer.ViewerType() == nextViewer.ViewerType() {
		return
	}

	nextViewer.SetActive(true)
	nextViewer.SetFocused(currViewer.IsFocused())
	nextViewer.SetFullScreen(currViewer.IsFullScreen())

	currViewer.SetActive(false)
	currViewer.SetFocused(false)
	currViewer.SetFullScreen(false)

	m.updateActiveViewer()
}

func (m *Model) updateActiveViewer() {
	if !m.ready {
		return
	}

	m.activeViewer().SetArticle(*m.activeArticle)
}

func (m *Model) updateActiveArticle() {
	newsIndex, articleIndex := m.selector.GetSelectedIndex()
	if newsIndex == 0 {
		m.activeArticle = &news.NationalNews[articleIndex]
	} else {
		m.activeArticle = &news.RegionalNews[articleIndex]
	}
}

func (m *Model) updateSizes(width, height int) {
	m.width = width
	m.height = height

	selectorWidthMultiplier := max(min(m.config.Settings.SelectorWidth, 0.8), 0.2)
	selectorWidth := int(float32(m.width) * selectorWidthMultiplier)

	m.selector.SetDims(selectorWidth, m.height-m.helper.Height()-5)

	isViewerFullscreen := m.activeViewer().IsFullScreen()
	for _, viewer := range m.viewers {
		if isViewerFullscreen {
			viewer.SetDims(m.width, m.height-m.helper.Height())
		} else {
			viewer.SetDims(m.width-selectorWidth-6, m.height-m.helper.Height())
		}
	}

	m.helper.SetWidth(m.width)
}

func (m Model) View() string {
	if m.loadingFailed {
		content := "Laden der Nachrichten fehlgeschlagen... press q to quit"
		return m.style.ScreenCenteredStyle(m.width, m.height).Render(content)
	}
	if !m.ready {
		content := fmt.Sprintf("%s Lade Nachrichten... press q to quit", m.spinner.View())
		return m.style.ScreenCenteredStyle(m.width, m.height).Render(content)
	}

	selector := m.selector.View()
	viewer := m.activeViewer().View()
	help := m.helper.View()

	return lipgloss.JoinHorizontal(lipgloss.Top, selector, viewer) + help
}
