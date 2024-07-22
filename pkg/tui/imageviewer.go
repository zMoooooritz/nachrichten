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
	cache *ImageCache
}

func NewImageViewer(viewer BaseViewer, ic *ImageCache) *ImageViewer {
	viewer.modeName = "Bild"
	return &ImageViewer{
		BaseViewer: viewer,
		image:      image.Rect(0, 0, 1, 1),
		cache:      ic,
	}
}

func (i ImageViewer) Update(msg tea.Msg) (Viewer, tea.Cmd) {
	var cmd tea.Cmd
	i.viewport, cmd = i.viewport.Update(msg)
	return &ImageViewer{BaseViewer: i.BaseViewer, image: i.image, cache: i.cache}, cmd
}

func (i *ImageViewer) SetArticle(article tagesschau.Article) {
	i.SetHeaderData(article)
	img := i.cache.GetImage(article.ID, article.ImageData.ImageVariants.RectSmall)
	i.pushImageToViewer(img)
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
