package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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
	repr := util.FormatParagraphs(paragraphs, r.viewport.Width-2, r.style)
	r.viewport.SetContent(repr)
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
