package tui

import (
	"fmt"
	"image"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/nachrichten/pkg/config"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
	"github.com/zMoooooritz/nachrichten/pkg/util"
)

type ImageViewer struct {
	Viewer
	image image.Image
}

func NewImageViewer(s config.Style, km viewport.KeyMap, isActive bool) ImageViewer {
	vp := viewport.New(0, 0)
	vp.KeyMap = km
	return ImageViewer{
		Viewer: Viewer{
			style:    s,
			isActive: isActive,
			viewport: vp,
		},
		image: image.Rect(0, 0, 1, 1),
	}
}

func (i *ImageViewer) SetArticle(article tagesschau.Article) {
	i.image = article.Thumbnail
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
