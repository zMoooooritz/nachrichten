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

type Viewer interface {
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
	Update(tea.Msg) (Viewer, tea.Cmd)
	View() string
}

type BaseViewer struct {
	viewerType   ViewerType
	style        config.Style
	isActive     bool
	isFocused    bool
	isFullScreen bool
	title        string
	date         string
	modeName     string
	viewport     viewport.Model
}

func NewViewer(viewerType ViewerType, s config.Style, km viewport.KeyMap, isActive bool) BaseViewer {
	vp := viewport.New(0, 0)
	vp.KeyMap = km
	return BaseViewer{
		viewerType: viewerType,
		style:      s,
		isActive:   isActive,
		viewport:   vp,
	}
}

func (v *BaseViewer) SetArticle(article tagesschau.Article) {
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

func (v *BaseViewer) SetHeaderData(title string, date time.Time) {
	v.title = title
	v.date = date.Format(germanDateFormat)
}

func (v BaseViewer) Init() tea.Cmd {
	return nil
}

func (v BaseViewer) Update(msg tea.Msg) (BaseViewer, tea.Cmd) {
	return v, nil
}

func (v BaseViewer) View() string {
	if !v.isActive {
		return ""
	}
	return fmt.Sprintf("%s\n%s\n%s", v.headerView(), v.viewport.View(), v.footerView())
}

func (v BaseViewer) headerView() string {
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

func (v BaseViewer) footerView() string {
	modeStyle := v.style.ReaderTitleInactiveStyle
	infoStyle := v.style.ReaderInfoInactiveStyle
	lineStyle := v.style.InactiveStyle
	fillCharacter := config.SingleFillCharacter
	if v.isFocused || v.isFullScreen {
		modeStyle = v.style.ReaderTitleActiveStyle
		infoStyle = v.style.ReaderInfoActiveStyle
		lineStyle = v.style.ActiveStyle
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
