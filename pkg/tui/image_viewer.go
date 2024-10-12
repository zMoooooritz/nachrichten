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
	viewer.modeName = "Bild"
	return &ImageViewer{
		BaseViewer: viewer,
		image:      image.Rect(0, 0, 1, 1),
	}
}

func (i *ImageViewer) Update(msg tea.Msg) (Viewer, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case UpdatedArticle:
		i.SetArticle(tagesschau.Article(msg))
	}

	if i.isFocused || i.isFullScreen {
		i.viewport, cmd = i.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	bv, cmd := i.BaseViewer.Update(msg)
	cmds = append(cmds, cmd)
	return &ImageViewer{BaseViewer: bv, image: i.image}, tea.Batch(cmds...)
}

func (i *ImageViewer) SetArticle(article tagesschau.Article) {
	i.SetHeaderData(article)
	if article.IsEmptyArticle() {
		i.viewport.SetContent("")
	} else {
		img := i.shared.imageCache.GetImage(article.ID, article.ImageData.ImageVariants.RectSmall)
		i.pushImageToViewer(img)
	}
}

func (i *ImageViewer) pushImageToViewer(img image.Image) {
	w := i.viewport.Width - 4
	h := i.viewport.Height - 2
	image := util.ImageToAscii(img, uint(w), uint(h), true)

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
