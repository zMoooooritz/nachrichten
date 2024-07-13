package tui

import (
	"fmt"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/zMoooooritz/nachrichten/pkg/config"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
	"github.com/zMoooooritz/nachrichten/pkg/util"
)

type Reader struct {
	Viewer
}

func NewReader(s config.Style, km viewport.KeyMap, isActive bool) Reader {
	vp := viewport.New(0, 0)
	vp.KeyMap = km
	return Reader{
		Viewer: Viewer{
			style:    s,
			isActive: isActive,
			viewport: vp,
		},
	}
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
