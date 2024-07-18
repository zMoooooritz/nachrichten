package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Details struct {
	BaseViewer
}

func NewDetails(viewer BaseViewer) *Details {
	viewer.modeName = "Details"
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
	details := ""
	if article.IsRegionalArticle() {
		details = d.buildRegionalArticleDetails(article)
	} else {
		details = d.buildNationalArticleDetails(article)
	}

	return d.style.InactiveStyle.Render(details)
}

func (d Details) buildRegionalArticleDetails(article tagesschau.Article) string {
	regionsStr := []string{}
	for _, regionID := range article.RegionIDs {
		regionsStr = append(regionsStr, string(tagesschau.GERMAN_NAMES[regionID]))
	}

	l := list.New()
	l.Item(fmt.Sprintf("Titel: %s", article.Desc))
	if len(regionsStr) == 1 {
		l.Item(fmt.Sprintf("Region: %s", regionsStr[0]))
	} else if len(regionsStr) > 0 {
		l.Item("Regionen:")
		l.Item(list.New(regionsStr))
	}
	caser := cases.Title(language.English)
	l.Item(fmt.Sprintf("Typ: %s", caser.String(article.Type)))
	l.Item(breakingText(article.Breaking))
	return l.String()
}

func (d Details) buildNationalArticleDetails(article tagesschau.Article) string {
	tagStr := []string{}
	for _, tag := range article.Tags {
		tagStr = append(tagStr, tag.Tag)
	}
	relatedStr := []string{}
	for index, related := range article.GetRelatedArticles() {
		relatedStr = append(relatedStr, fmt.Sprintf("[%d]\n   %s\n    %s", index+1, strings.TrimSpace(related.Topline), strings.TrimSpace(related.Desc)))
	}
	l := list.New()
	l.Item(fmt.Sprintf("Titel: %s", article.Topline))
	l.Item(fmt.Sprintf("Untertitel: %s", article.Desc))
	caser := cases.Title(language.English)
	if article.IsRegionalArticle() {
		l.Item(fmt.Sprintf("Region: %s", tagesschau.GERMAN_NAMES[article.RegionID]))
	} else {
		l.Item(fmt.Sprintf("Ressort: %s", caser.String(article.Ressort)))
	}
	l.Item(fmt.Sprintf("Typ: %s", caser.String(article.Type)))
	l.Item(breakingText(article.Breaking))
	l.Item("Tags:")
	l.Item(list.New(tagStr))
	if len(relatedStr) > 0 {
		l.Item("Verwandt:")
		l.Item(list.New(relatedStr))
	}
	return d.style.InactiveStyle.Render(l.String())
}

func breakingText(breaking bool) string {
	if breaking {
		return "Eilmeldung: Ja"
	} else {
		return "Eilmeldung: Nein"
	}
}
