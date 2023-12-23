package tui

import (
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/nachrichten/pkg/config"
	"github.com/zMoooooritz/nachrichten/pkg/util"
)

type ImageViewer struct {
	style        config.Style
	isActive     bool
	isFocused    bool
	isFullScreen bool
	toplineText  string
	dateText     string
	viewport     viewport.Model
	image        image.Image
}

func NewImageViewer(s config.Style) ImageViewer {
	return ImageViewer{
		style:    s,
		isActive: false,
		viewport: viewport.New(0, 0),
		image:    image.Rect(0, 0, 1, 1),
	}
}

func (i *ImageViewer) SetActive(isActive bool) {
	i.isActive = isActive
}

func (i *ImageViewer) IsActive() bool {
	return i.isActive
}

func (i *ImageViewer) SetFocused(isFocused bool) {
	i.isFocused = isFocused
}

func (i *ImageViewer) IsFocused() bool {
	return i.isFocused
}

func (i *ImageViewer) SetFullScreen(isFullScreen bool) {
	i.isFullScreen = isFullScreen
}

func (i *ImageViewer) IsFullScreen() bool {
	return i.isFullScreen
}

func (i *ImageViewer) SetDims(w, h int) {
	i.viewport.Width = w
	i.viewport.Height = h - lipgloss.Height(i.headerView()) - lipgloss.Height(i.footerView())
	i.viewport.YPosition = lipgloss.Height(i.headerView())
	i.PushImageToViewer()
}

func (i *ImageViewer) SetImage(img image.Image) {
	i.image = img
	i.PushImageToViewer()
}

func (i *ImageViewer) PushImageToViewer() {
	w := i.viewport.Width - 4
	h := i.viewport.Height - 2
	image := util.ImageToAscii(i.image, uint(w), uint(h), true)

	strRepr := ""
	for _, row := range image {
		rowRepr := ""
		for _, char := range row {
			rowRepr += char
		}
		strRepr += lipgloss.PlaceHorizontal(i.viewport.Width, lipgloss.Center, rowRepr) + "\n"
	}

	strRepr = lipgloss.PlaceVertical(h, lipgloss.Center, strRepr)
	i.viewport.SetContent(strRepr)
}

func (i *ImageViewer) SetHeaderContent(topline string, date time.Time) {
	i.toplineText = topline
	i.dateText = date.Format(germanDateFormat)
}

func (i ImageViewer) Init() tea.Cmd {
	return nil
}

func (i ImageViewer) Update(msg tea.Msg) (ImageViewer, tea.Cmd) {
	var cmd tea.Cmd
	i.viewport, cmd = i.viewport.Update(msg)
	return i, tea.Batch(cmd)
}

func (i ImageViewer) View() string {
	if !i.isActive {
		return ""
	}
	return fmt.Sprintf("%s\n%s\n%s", i.headerView(), i.viewport.View(), i.footerView())
}

func (i ImageViewer) headerView() string {
	titleStyle := i.style.ReaderTitleInactiveStyle
	lineStyle := i.style.InactiveStyle
	dateStyle := i.style.ReaderInfoInactiveStyle
	fillCharacter := config.SingleFillCharacter
	if i.isFocused || i.isFullScreen {
		titleStyle = i.style.ReaderTitleActiveStyle
		lineStyle = i.style.ActiveStyle
		dateStyle = i.style.ReaderInfoActiveStyle
		fillCharacter = config.DoubleFillCharacter
	}

	title := titleStyle.Render(i.toplineText)
	date := dateStyle.Render(i.dateText)
	line := lineStyle.Render(strings.Repeat(fillCharacter, util.Max(0, i.viewport.Width-lipgloss.Width(title)-lipgloss.Width(date))))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line, date)
}

func (i ImageViewer) footerView() string {
	infoStyle := i.style.ReaderInfoInactiveStyle
	lineStyle := i.style.InactiveStyle
	fillCharacter := config.SingleFillCharacter
	if i.isFocused || i.isFullScreen {
		infoStyle = i.style.ReaderInfoActiveStyle
		lineStyle = i.style.ActiveStyle
		fillCharacter = config.DoubleFillCharacter
	}

	info := infoStyle.Render(fmt.Sprintf("%3.f%%", i.viewport.ScrollPercent()*100))
	line := lineStyle.Render(strings.Repeat(fillCharacter, util.Max(0, i.viewport.Width-lipgloss.Width(info))))

	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
