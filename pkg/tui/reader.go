package tui

import (
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
	return &Reader{
		BaseViewer: viewer,
	}
}

func (r Reader) Update(msg tea.Msg) (Viewer, tea.Cmd) {
	var cmd tea.Cmd
	r.viewport, cmd = r.viewport.Update(msg)
	return &Reader{BaseViewer: r.BaseViewer}, cmd
}

func (r *Reader) SetArticle(article tagesschau.Article) {
	paragraphs := tagesschau.ContentToParagraphs(article.Content)
	r.viewport.SetContent(r.formatParagraphs(paragraphs))
}

func (r Reader) formatParagraphs(paragraphs []string) string {
	width := r.viewport.Width - 2
	options := md.Options{
		EscapeMode: "disabled",
	}
	converter := md.NewConverter("", true, &options)
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
		glamour.WithStyles(r.style.ReaderStyle),
	)

	result := ""
	for _, p := range paragraphs {
		text, _ := converter.ConvertString(p)
		text, _ = renderer.Render(text)
		result += text
	}
	return util.PadText(result, width)
}
