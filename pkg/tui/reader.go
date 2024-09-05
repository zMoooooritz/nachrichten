package tui

import (
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
	"github.com/zMoooooritz/nachrichten/pkg/util"
)

type Reader struct {
	BaseViewer
}

func NewReader(viewer BaseViewer) *Reader {
	viewer.modeName = "Artikel"
	return &Reader{
		BaseViewer: viewer,
	}
}

func (r Reader) Update(msg tea.Msg) (Viewer, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	if r.IsFocused() || r.isFullScreen {
		r.viewport, cmd = r.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	bv, cmd := r.BaseViewer.Update(msg)
	cmds = append(cmds, cmd)
	return &Reader{BaseViewer: bv}, tea.Batch(cmds...)
}

func (r *Reader) SetArticle(article tagesschau.Article) {
	r.SetHeaderData(article)
	paragraphs := tagesschau.ContentToParagraphs(article.Content)
	r.viewport.SetContent(r.formatParagraphs(paragraphs))
}

func (r Reader) formatParagraphs(paragraphs []string) string {
	width := r.viewport.Width - 6
	options := md.Options{
		EscapeMode: "disabled",
	}
	converter := md.NewConverter("", true, &options)
	renderer, err := glamour.NewTermRenderer(
		glamour.WithWordWrap(width),
		glamour.WithStyles(r.style.ReaderStyle),
	)
	if err != nil {
		util.Logger.Fatalln(err)
		return util.PadText("Unable to parse and print article", width)
	}

	joined := strings.Join(paragraphs, "\n\n")
	text, err := converter.ConvertString(joined)
	if err != nil {
		util.Logger.Fatalln(err)
		return util.PadText("Unable to parse and print article", width)
	}
	result, err := renderer.Render(text)
	if err != nil {
		util.Logger.Fatalln(err)
		return util.PadText("Unable to parse and print article", width)
	}
	return util.PadText(result, width)
}
