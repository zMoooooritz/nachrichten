package tui

import (
	"fmt"
	"strings"
	"time"

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

type ViewerImplementation interface {
	SetArticle(tagesschau.Article)
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
	SetHeaderData(string, time.Time)
	Init() tea.Cmd
	Update(tea.Msg) (ViewerImplementation, tea.Cmd)
	View() string
}

type Viewer struct {
	viewerType   ViewerType
	style        config.Style
	isActive     bool
	isFocused    bool
	isFullScreen bool
	title        string
	date         string
	viewport     viewport.Model
}

func NewViewer(viewerType ViewerType, s config.Style, km viewport.KeyMap, isActive bool) Viewer {
	vp := viewport.New(0, 0)
	vp.KeyMap = km
	return Viewer{
		viewerType: viewerType,
		style:      s,
		isActive:   isActive,
		viewport:   vp,
	}
}

func (v *Viewer) SetArticle(article tagesschau.Article) {
}

func (v Viewer) ViewerType() ViewerType {
	return v.viewerType
}

func (v *Viewer) GotoTop() {
	v.viewport.GotoTop()
}

func (v *Viewer) GotoBottom() {
	v.viewport.GotoBottom()
}

func (v *Viewer) SetActive(isActive bool) {
	v.isActive = isActive
}

func (v *Viewer) IsActive() bool {
	return v.isActive
}

func (v *Viewer) SetFocused(isFocused bool) {
	v.isFocused = isFocused
}

func (v *Viewer) IsFocused() bool {
	return v.isFocused
}

func (v *Viewer) SetFullScreen(isFullScreen bool) {
	v.isFullScreen = isFullScreen
}

func (v *Viewer) IsFullScreen() bool {
	return v.isFullScreen
}

func (v *Viewer) SetDims(w, h int) {
	v.viewport.Width = w
	v.viewport.Height = h - lipgloss.Height(v.headerView()) - lipgloss.Height(v.footerView())
	v.viewport.YPosition = lipgloss.Height(v.headerView())
}

func (v *Viewer) SetHeaderData(title string, date time.Time) {
	v.title = title
	v.date = date.Format(germanDateFormat)
}

func (v Viewer) Init() tea.Cmd {
	return nil
}

func (v Viewer) Update(msg tea.Msg) (Viewer, tea.Cmd) {
	return v, nil
}

func (v Viewer) View() string {
	if !v.isActive {
		return ""
	}
	return fmt.Sprintf("%s\n%s\n%s", v.headerView(), v.viewport.View(), v.footerView())
}

func (v Viewer) headerView() string {
	titleStyle := v.style.ReaderTitleInactiveStyle
	lineStyle := v.style.InactiveStyle
	dateStyle := v.style.ReaderInfoInactiveStyle
	fillCharacter := config.SingleFillCharacter
	if v.isFocused || v.isFullScreen {
		titleStyle = v.style.ReaderTitleActiveStyle
		lineStyle = v.style.ActiveStyle
		dateStyle = v.style.ReaderInfoActiveStyle
		fillCharacter = config.DoubleFillCharacter
	}

	title := titleStyle.Render(v.title)
	date := dateStyle.Render(v.date)
	line := lineStyle.Render(strings.Repeat(fillCharacter, max(0, v.viewport.Width-lipgloss.Width(title)-lipgloss.Width(date))))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line, date)
}

func (v Viewer) footerView() string {
	infoStyle := v.style.ReaderInfoInactiveStyle
	lineStyle := v.style.InactiveStyle
	fillCharacter := config.SingleFillCharacter
	if v.isFocused || v.isFullScreen {
		infoStyle = v.style.ReaderInfoActiveStyle
		lineStyle = v.style.ActiveStyle
		fillCharacter = config.DoubleFillCharacter
	}

	info := infoStyle.Render(fmt.Sprintf("%3.f%%", v.viewport.ScrollPercent()*100))
	line := lineStyle.Render(strings.Repeat(fillCharacter, max(0, v.viewport.Width-lipgloss.Width(info))))

	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
