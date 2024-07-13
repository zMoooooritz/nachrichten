package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
)

type Details struct {
	BaseViewer
}

func NewDetails(viewer BaseViewer) *Details {
	return &Details{
		BaseViewer: viewer,
	}
}

func (d Details) Update(msg tea.Msg) (Viewer, tea.Cmd) {
	var cmd tea.Cmd
	d.viewport, cmd = d.viewport.Update(msg)
	return &Details{BaseViewer: d.BaseViewer}, cmd
}

func (d *Details) SetArticle(article tagesschau.Article) {
	d.viewport.SetContent(d.buildDetails(article))
}

func (d Details) buildDetails(article tagesschau.Article) string {
	tagStr := []string{}
	for _, tag := range article.Tags {
		tagStr = append(tagStr, tag.Tag)
	}
	title := fmt.Sprintf("Titel: %s", article.Desc)
	region := fmt.Sprintf("Region: %s", tagesschau.GERMAN_NAMES[article.RegionID])
	ressort := fmt.Sprintf("Ressort: %s", strings.Title(article.Ressort))
	breaking := ""
	if article.Breaking {
		breaking = fmt.Sprintf("Eilmeldung: ja")
	} else {
		breaking = fmt.Sprintf("Eilmeldung: nein")
	}
	l := list.New(
		title, list.New(),
		region, list.New(),
		ressort, list.New(),
		"Tags", list.New(tagStr),
		breaking, list.New(),
	)
	return l.String()
}
