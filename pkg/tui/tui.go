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
	news tagesschau.News
)

type Model struct {
	opener        util.Opener
	ready         bool
	loadingFailed bool
	shared        *SharedState
	navigator     *Navigator
	viewManager   *ViewManager
	helper        *Helper
	spinner       spinner.Model
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
	imageCache    *ImageCache
}

func InitialModel(c config.Configuration) Model {
	initialHelpState := HS_NORMAL
	if c.Settings.HideHelpOnStartup {
		initialHelpState = HS_HIDDEN
	}

	style := config.NewsStyle(c.Theme)
	shared := &SharedState{
		mode:       NORMAL_MODE,
		style:      style,
		keys:       c.Keys,
		keymap:     GetKeyMap(c.Keys),
		config:     c,
		imageCache: NewImageCache(),
	}

	return Model{
		opener:      util.NewOpener(c.Applications),
		ready:       false,
		helper:      NewHelper(shared, initialHelpState),
		navigator:   NewNavigator(shared),
		shared:      shared,
		viewManager: NewViewManager(shared),
		spinner:     NewDotSpinner(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(loadNews, m.spinner.Tick)
}

func refreshFunc(article tagesschau.Article) tea.Cmd {
	return func() tea.Msg {
		return UpdatedArticle(article)
	}
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
			go m.shared.imageCache.LoadThumbnails(append(news.NationalNews, news.RegionalNews...))
		}
		m.ready = true
		m.shared.activeArticle = news.NationalNews[0]
		cmds = append(cmds, refreshFunc(m.shared.activeArticle))
	case UpdatedArticle:
		article := tagesschau.Article(msg)
		if m.shared.config.Settings.PreloadThumbnails {
			go m.shared.imageCache.LoadThumbnail(article)
		}
		m.shared.activeArticle = article
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
		m.width = msg.Width
		m.height = msg.Height
		cmds = append(cmds, refreshFunc(m.shared.activeArticle))
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

	m.viewManager, cmd = m.viewManager.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.shared.keymap.full):
			cmds = append(cmds, refreshFunc(m.shared.activeArticle))
		}
	}

	return m, tea.Batch(cmds...)
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

	m.helper.SetWidth(m.width)
	help := m.helper.View()

	navigatorWidthMultiplier := max(min(m.shared.config.Settings.NavigatorWidth, 0.8), 0.2)
	navigatorWidth := int(float32(m.width) * navigatorWidthMultiplier)

	helperHeight := lipgloss.Height(help)
	if !m.helper.IsVisible() {
		helperHeight = 0
	}

	m.navigator.SetDims(navigatorWidth, m.height-helperHeight)
	navigator := m.navigator.View()

	m.viewManager.SetDims(m.width, m.height-helperHeight, lipgloss.Width(navigator))
	viewer := m.viewManager.View()

	view := lipgloss.JoinHorizontal(lipgloss.Top, navigator, viewer)
	if m.helper.IsVisible() {
		view = lipgloss.JoinVertical(lipgloss.Center, view, help)
	}
	return view
}
