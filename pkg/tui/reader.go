package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/nachrichten/pkg/config"
	"github.com/zMoooooritz/nachrichten/pkg/util"
)

type Reader struct {
	style        config.Style
	isActive     bool
	isFocused    bool
	isFullScreen bool
	toplineText  string
	dateText     string
	viewport     viewport.Model
}

func NewReader(s config.Style, km viewport.KeyMap) Reader {
	vp := viewport.New(0, 0)
	vp.KeyMap = km
	return Reader{
		style:    s,
		isActive: true,
		viewport: vp,
	}
}

func (r *Reader) GotoTop() {
	r.viewport.GotoTop()
}

func (r *Reader) GotoBottom() {
	r.viewport.GotoBottom()
}

func (r *Reader) SetActive(isActive bool) {
	r.isActive = isActive
}

func (r *Reader) IsActive() bool {
	return r.isActive
}

func (r *Reader) SetFocused(isFocused bool) {
	r.isFocused = isFocused
}

func (r *Reader) IsFocused() bool {
	return r.isFocused
}

func (r *Reader) SetFullScreen(isFullScreen bool) {
	r.isFullScreen = isFullScreen
}

func (r *Reader) IsFullScreen() bool {
	return r.isFullScreen
}

func (r *Reader) SetDims(w, h int) {
	r.viewport.Width = w
	r.viewport.Height = h - lipgloss.Height(r.headerView()) - lipgloss.Height(r.footerView())
	r.viewport.YPosition = lipgloss.Height(r.headerView())
}

func (r *Reader) SetContent(paragraphs []string) {
	repr := util.FormatParagraphs(paragraphs, r.viewport.Width, r.style)
	r.viewport.SetContent(repr)
}

func (r *Reader) SetHeaderContent(topline string, date time.Time) {
	r.toplineText = topline
	r.dateText = date.Format(germanDateFormat)
}

func (r Reader) Init() tea.Cmd {
	return nil
}

func (r Reader) Update(msg tea.Msg) (Reader, tea.Cmd) {
	var cmd tea.Cmd
	r.viewport, cmd = r.viewport.Update(msg)
	return r, tea.Batch(cmd)
}

func (r Reader) View() string {
	if !r.isActive {
		return ""
	}
	return fmt.Sprintf("%s\n%s\n%s", r.headerView(), r.viewport.View(), r.footerView())
}

func (r Reader) headerView() string {
	titleStyle := r.style.ReaderTitleInactiveStyle
	lineStyle := r.style.InactiveStyle
	dateStyle := r.style.ReaderInfoInactiveStyle
	fillCharacter := config.SingleFillCharacter
	if r.isFocused || r.isFullScreen {
		titleStyle = r.style.ReaderTitleActiveStyle
		lineStyle = r.style.ActiveStyle
		dateStyle = r.style.ReaderInfoActiveStyle
		fillCharacter = config.DoubleFillCharacter
	}

	title := titleStyle.Render(r.toplineText)
	date := dateStyle.Render(r.dateText)
	line := lineStyle.Render(strings.Repeat(fillCharacter, util.Max(0, r.viewport.Width-lipgloss.Width(title)-lipgloss.Width(date))))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line, date)
}

func (r Reader) footerView() string {
	infoStyle := r.style.ReaderInfoInactiveStyle
	lineStyle := r.style.InactiveStyle
	fillCharacter := config.SingleFillCharacter
	if r.isFocused || r.isFullScreen {
		infoStyle = r.style.ReaderInfoActiveStyle
		lineStyle = r.style.ActiveStyle
		fillCharacter = config.DoubleFillCharacter
	}

	info := infoStyle.Render(fmt.Sprintf("%3.f%%", r.viewport.ScrollPercent()*100))
	line := lineStyle.Render(strings.Repeat(fillCharacter, util.Max(0, r.viewport.Width-lipgloss.Width(info))))

	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
