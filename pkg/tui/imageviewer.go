package tui

import (
	"image"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
	"github.com/zMoooooritz/nachrichten/pkg/util"
)

type ImageViewer struct {
	BaseViewer
	image image.Image
}

func NewImageViewer(viewer BaseViewer) *ImageViewer {
	return &ImageViewer{
		BaseViewer: viewer,
		image:      image.Rect(0, 0, 1, 1),
	}
}

func (i ImageViewer) Update(msg tea.Msg) (Viewer, tea.Cmd) {
	var cmd tea.Cmd
	i.viewport, cmd = i.viewport.Update(msg)
	return &ImageViewer{BaseViewer: i.BaseViewer, image: i.image}, cmd
}

func (i *ImageViewer) SetArticle(article tagesschau.Article) {
	i.image = article.Thumbnail
	i.pushImageToViewer()
}

func (i *ImageViewer) pushImageToViewer() {
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
