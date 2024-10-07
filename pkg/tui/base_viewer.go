package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/nachrichten/pkg/config"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
)

type ViewerType int

const (
	VT_TEXT ViewerType = iota
	VT_IMAGE
	VT_DETAILS
)

type Viewer interface {
	SetArticle(tagesschau.Article)
	SetHeaderData(tagesschau.Article)
	ViewerType() ViewerType
	GotoTop()
	GotoBottom()
	SetActive(bool)
	IsActive() bool
	SetFocused(bool)
	IsFocused() bool
	SetFullScreen(bool)
	IsFullScreen() bool
	SetDims(int, int)
	Init() tea.Cmd
	Update(tea.Msg) (Viewer, tea.Cmd)
	View() string
}

type BaseViewer struct {
	viewerType   ViewerType
	shared       *SharedState
	isActive     bool
	isFocused    bool
	isFullScreen bool
	title        string
	date         string
	modeName     string
	viewport     viewport.Model
}

func NewViewer(viewerType ViewerType, shared *SharedState, isActive bool) BaseViewer {
	vp := viewport.New(0, 0)
	vp.KeyMap = ViewportKeymap(shared.keys)
	return BaseViewer{
		shared:     shared,
		viewerType: viewerType,
		isActive:   isActive,
		viewport:   vp,
	}
}

func (v *BaseViewer) SetArticle(article tagesschau.Article) {
	v.SetHeaderData(article)
}

func (v BaseViewer) GetModeText() string {
	return ""
}

func (v BaseViewer) ViewerType() ViewerType {
	return v.viewerType
}

func (v *BaseViewer) GotoTop() {
	v.viewport.GotoTop()
}

func (v *BaseViewer) GotoBottom() {
	v.viewport.GotoBottom()
}

func (v *BaseViewer) SetActive(isActive bool) {
	v.isActive = isActive
}

func (v *BaseViewer) IsActive() bool {
	return v.isActive
}

func (v *BaseViewer) SetFocused(isFocused bool) {
	v.isFocused = isFocused
}

func (v *BaseViewer) IsFocused() bool {
	return v.isFocused
}

func (v *BaseViewer) SetFullScreen(isFullScreen bool) {
	v.isFullScreen = isFullScreen
}

func (v *BaseViewer) IsFullScreen() bool {
	return v.isFullScreen
}

func (v *BaseViewer) SetDims(w, h int) {
	v.viewport.Width = w
	v.viewport.Height = h - lipgloss.Height(v.headerView()) - lipgloss.Height(v.footerView())
	v.viewport.YPosition = lipgloss.Height(v.headerView())
}

func (v *BaseViewer) SetHeaderData(article tagesschau.Article) {
	if article.IsRegionalArticle() {
		v.title = article.Desc
	} else {
		if article.Topline != "" {
			v.title = article.Topline
		} else {
			v.title = article.Desc
		}
	}
	v.date = article.Date.Format(germanDateFormat)
}

func (v BaseViewer) Init() tea.Cmd {
	return nil
}

func (v BaseViewer) Update(msg tea.Msg) (BaseViewer, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if v.shared.mode == INSERT_MODE {
			break
		}

		switch {
		case key.Matches(msg, v.shared.keymap.right):
			if v.isActive {
				v.isFocused = true
			}
		case key.Matches(msg, v.shared.keymap.left):
			if v.isActive {
				v.isFocused = false
			}
		case key.Matches(msg, v.shared.keymap.start):
			if v.isFocused || v.isFullScreen {
				v.GotoTop()
			}
		case key.Matches(msg, v.shared.keymap.end):
			if v.isFocused || v.isFullScreen {
				v.GotoBottom()
			}
		case key.Matches(msg, v.shared.keymap.full):
			if v.isActive {
				v.isFullScreen = !v.isFullScreen
			}
		}
	}

	return v, nil
}

func (v BaseViewer) View() string {
	if !v.isActive {
		return ""
	}
	return fmt.Sprintf("%s\n%s\n%s", v.headerView(), v.viewport.View(), v.footerView())
}

func (v BaseViewer) headerView() string {
	titleStyle := v.shared.style.ReaderTitleInactiveStyle
	lineStyle := v.shared.style.InactiveStyle
	dateStyle := v.shared.style.ReaderInfoInactiveStyle
	fillCharacter := config.SingleFillCharacter
	if v.isFocused || v.isFullScreen {
		titleStyle = v.shared.style.ReaderTitleActiveStyle
		lineStyle = v.shared.style.ActiveStyle
		dateStyle = v.shared.style.ReaderInfoActiveStyle
		fillCharacter = config.DoubleFillCharacter
	}

	title := titleStyle.Render(v.title)
	date := dateStyle.Render(v.date)
	line := lineStyle.Render(strings.Repeat(fillCharacter, max(0, v.viewport.Width-lipgloss.Width(title)-lipgloss.Width(date))))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line, date)
}

func (v BaseViewer) footerView() string {
	modeStyle := v.shared.style.ReaderTitleInactiveStyle
	infoStyle := v.shared.style.ReaderInfoInactiveStyle
	lineStyle := v.shared.style.InactiveStyle
	fillCharacter := config.SingleFillCharacter
	if v.isFocused || v.isFullScreen {
		modeStyle = v.shared.style.ReaderTitleActiveStyle
		infoStyle = v.shared.style.ReaderInfoActiveStyle
		lineStyle = v.shared.style.ActiveStyle
		fillCharacter = config.DoubleFillCharacter
	}

	mode := ""
	if v.modeName != "" {
		mode = modeStyle.Render(v.modeName)
	}
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", v.viewport.ScrollPercent()*100))
	line := lineStyle.Render(strings.Repeat(fillCharacter, max(0, v.viewport.Width-lipgloss.Width(mode)-lipgloss.Width(info))))

	return lipgloss.JoinHorizontal(lipgloss.Center, mode, line, info)
}
