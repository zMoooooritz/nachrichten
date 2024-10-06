package tui

import (
	"fmt"

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
	news        tagesschau.News
	refreshFunc = func() tea.Msg { return RefreshActiveViewer{} }
)

type Model struct {
	opener        util.Opener
	ready         bool
	loadingFailed bool
	shared        *SharedState
	helper        *Helper
	navigator     *Navigator
	viewers       []Viewer
	spinner       spinner.Model
	imageCache    *ImageCache
	width         int
	height        int
}

type Mode int

const (
	NORMAL_MODE Mode = iota
	INSERT_MODE
)

type SharedState struct {
	mode          Mode
	style         config.Style
	keys          config.Keys
	keymap        KeyMap
	config        config.Configuration
	activeArticle tagesschau.Article
}

func InitialModel(c config.Configuration) Model {
	initialHelpState := HS_NORMAL
	if c.Settings.HideHelpOnStartup {
		initialHelpState = HS_HIDDEN
	}

	ic := NewImageCache()

	style := config.NewsStyle(c.Theme)
	shared := &SharedState{
		mode:   NORMAL_MODE,
		style:  style,
		keys:   c.Keys,
		keymap: GetKeyMap(c.Keys),
		config: c,
	}

	viewers := []Viewer{
		NewReader(NewViewer(VT_TEXT, shared, true)),
		NewImageViewer(NewViewer(VT_IMAGE, shared, false), ic),
		NewDetails(NewViewer(VT_DETAILS, shared, false)),
	}

	return Model{
		opener:     util.NewOpener(c.Applications),
		ready:      false,
		helper:     NewHelper(shared, initialHelpState),
		navigator:  NewNavigator(shared),
		shared:     shared,
		viewers:    viewers,
		spinner:    NewDotSpinner(),
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
		if m.shared.config.Settings.PreloadThumbnails {
			go m.preloadMainThumbnails()
		}
		m.ready = true
		m.shared.activeArticle = news.NationalNews[0]
		cmds = append(cmds, refreshFunc)
	case ChangedActiveArticle:
		article := tagesschau.Article(msg)
		if m.shared.config.Settings.PreloadThumbnails {
			go m.preloadThumbnail(article)
		}
		m.shared.activeArticle = article
	case ShowTextViewer:
		m.showViewer(VT_TEXT)
	case tea.KeyMsg:
		if m.shared.mode != NORMAL_MODE {
			break
		}

		if key.Matches(msg, m.shared.keymap.quit) {
			return m, tea.Quit
		}

		if !m.ready {
			return m, nil
		}

		switch {
		case key.Matches(msg, m.shared.keymap.article):
			m.showViewer(VT_TEXT)
		case key.Matches(msg, m.shared.keymap.image):
			m.showViewer(VT_IMAGE)
		case key.Matches(msg, m.shared.keymap.details):
			m.showViewer(VT_DETAILS)
		case key.Matches(msg, m.shared.keymap.open):
			m.opener.OpenUrl(util.TypeHTML, m.shared.activeArticle.URL)
		case key.Matches(msg, m.shared.keymap.video):
			m.opener.OpenUrl(util.TypeVideo, m.shared.activeArticle.Video.VideoVariants.Big)
		case key.Matches(msg, m.shared.keymap.shortNews):
			url, err := tagesschau.GetShortNewsURL()
			if err == nil {
				m.opener.OpenUrl(util.TypeVideo, url)
			}
		}
	case tea.WindowSizeMsg:
		m.updateSizes(msg.Width, msg.Height)
		cmds = append(cmds, refreshFunc)
	}

	if !m.ready {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	m.helper, cmd = m.helper.Update(msg)
	cmds = append(cmds, cmd)

	m.navigator, cmd = m.navigator.Update(msg)
	cmds = append(cmds, cmd)

	for i, viewer := range m.viewers {
		updatedViewer, cmd := viewer.Update(msg)
		m.viewers[i] = updatedViewer
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.shared.keymap.help):
			m.updateSizes(m.width, m.height)
		case key.Matches(msg, m.shared.keymap.full):
			m.updateSizes(m.width, m.height)
			cmds = append(cmds, refreshFunc)
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

func (m *Model) preloadMainThumbnails() {
	imageSpec := tagesschau.ImageSpec{Size: tagesschau.SMALL, Ratio: tagesschau.RECT}
	for _, a := range news.NationalNews {
		_ = m.imageCache.LoadImage(a.ID, tagesschau.GetImageURL(a.ImageData.ImageVariants, imageSpec))
	}
	for _, a := range news.RegionalNews {
		_ = m.imageCache.LoadImage(a.ID, tagesschau.GetImageURL(a.ImageData.ImageVariants, imageSpec))
	}
}

func (m *Model) preloadThumbnail(article tagesschau.Article) {
	imageSpec := tagesschau.ImageSpec{Size: tagesschau.SMALL, Ratio: tagesschau.RECT}
	_ = m.imageCache.LoadImage(article.ID, tagesschau.GetImageURL(article.ImageData.ImageVariants, imageSpec))
}

func (m *Model) showViewer(vt ViewerType) tea.Cmd {
	currViewer := m.activeViewer()
	nextViewer := m.activeViewer()
	for _, viewer := range m.viewers {
		if viewer.ViewerType() == vt {
			nextViewer = viewer
		}
	}
	if currViewer.ViewerType() == nextViewer.ViewerType() {
		return nil
	}

	nextViewer.SetActive(true)
	nextViewer.SetFocused(currViewer.IsFocused())
	nextViewer.SetFullScreen(currViewer.IsFullScreen())

	currViewer.SetActive(false)
	currViewer.SetFocused(false)
	currViewer.SetFullScreen(false)

	return refreshFunc
}

func (m *Model) updateSizes(width, height int) {
	m.width = width
	m.height = height

	navigatorWidthMultiplier := max(min(m.shared.config.Settings.NavigatorWidth, 0.8), 0.2)
	navigatorWidth := int(float32(m.width) * navigatorWidthMultiplier)

	m.navigator.SetDims(navigatorWidth, m.height-m.helper.Height()-5)

	isViewerFullscreen := m.activeViewer().IsFullScreen()
	for _, viewer := range m.viewers {
		if isViewerFullscreen {
			viewer.SetDims(m.width, m.height-m.helper.Height())
		} else {
			viewer.SetDims(m.width-navigatorWidth-6, m.height-m.helper.Height())
		}
	}

	m.helper.SetWidth(m.width)
}

func (m Model) View() string {
	if m.loadingFailed {
		content := "Laden der Nachrichten fehlgeschlagen... press q to quit"
		return m.shared.style.ScreenCenteredStyle(m.width, m.height).Render(content)
	}
	if !m.ready {
		content := fmt.Sprintf("%s Lade Nachrichten... press q to quit", m.spinner.View())
		return m.shared.style.ScreenCenteredStyle(m.width, m.height).Render(content)
	}

	navigator := m.navigator.View()
	viewer := m.activeViewer().View()
	help := m.helper.View()

	return lipgloss.JoinHorizontal(lipgloss.Top, navigator, viewer) + help
}
