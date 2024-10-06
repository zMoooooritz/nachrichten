package tui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func (d *Details) Update(msg tea.Msg) (Viewer, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case ChangedActiveArticle:
		d.SetArticle(tagesschau.Article(msg))
	case RefreshActiveViewer:
		d.SetArticle(d.shared.activeArticle)
	}

	if d.isFocused || d.isFullScreen {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			keyStr := msg.String()
			if keyStr >= "0" && keyStr <= "9" {
				keyInt, _ := strconv.Atoi(keyStr)
				cmd = d.handleNumberInput(keyInt)
				cmds = append(cmds, cmd)
			}
		}

		d.viewport, cmd = d.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	bv, cmd := d.BaseViewer.Update(msg)
	cmds = append(cmds, cmd)
	return &Details{BaseViewer: bv}, tea.Batch(cmds...)
}

func (d *Details) handleNumberInput(number int) tea.Cmd {
	related := d.shared.activeArticle.GetRelatedArticles()
	index := number - 1
	if 0 <= index && index < len(related) {
		article, err := tagesschau.LoadArticle(related[index].Details)
		if err == nil {
			return tea.Batch(
				func() tea.Msg {
					return ChangedActiveArticle(*article)
				},
				func() tea.Msg {
					return ShowTextViewer{}
				},
			)
		}
	}
	return nil
}

func (d *Details) SetArticle(article tagesschau.Article) {
	d.SetHeaderData(article)
	d.viewport.SetContent(d.buildDetails(article))
}

func (d Details) buildDetails(article tagesschau.Article) string {
	details := ""
	if article.IsRegionalArticle() {
		details = d.buildRegionalArticleDetails(article)
	} else {
		details = d.buildNationalArticleDetails(article)
	}

	return d.shared.style.InactiveStyle.Render(details)
}

func (d Details) buildRegionalArticleDetails(article tagesschau.Article) string {
	regionStr := ""
	regionsStr := ""
	for _, regionID := range article.RegionIDs {
		name, err := tagesschau.RegionIdToName(int(regionID))
		if err == nil {
			if regionStr == "" {
				regionStr = name
			}
			regionsStr += "  - " + name + "\n"
		}
	}

	strBuf := d.renderEntry("Titel", article.Desc) + "\n"

	if len(article.RegionIDs) == 1 {
		strBuf += d.renderEntry("Region", regionStr) + "\n"
	} else {
		strBuf += d.shared.style.ActiveHighlightStyle.Render("Regionen:") + "\n"
		strBuf += regionsStr
	}
	caser := cases.Title(language.English)
	strBuf += d.renderEntry("Typ", caser.String(article.Type)) + "\n"
	strBuf += d.breakingText(article.Breaking) + "\n"

	return lipgloss.NewStyle().PaddingLeft(2).Render(strBuf)
}

func (d Details) buildNationalArticleDetails(article tagesschau.Article) string {
	tagStr := ""
	for _, tag := range article.Tags {
		tagStr += "  - " + tag.Tag + "\n"
	}
	relatedStr := ""
	for index, related := range article.GetRelatedArticles() {
		ident := d.shared.style.HighlightStyle.Render(fmt.Sprintf("  [%d] ", index+1))
		repr := strings.TrimSpace(related.Topline)
		if strings.TrimSpace(related.Desc) != "" {
			repr += " - " + strings.TrimSpace(related.Desc)
		}
		relatedStr += ident + d.shared.style.InactiveStyle.Render(repr) + "\n"
	}

	strBuf := d.renderEntry("Titel", article.Topline) + "\n"
	strBuf += d.renderEntry("Untertitel", article.Desc) + "\n"
	caser := cases.Title(language.English)
	if article.Ressort != "" {
		strBuf += d.renderEntry("Ressort", caser.String(article.Ressort)) + "\n"
	}
	strBuf += d.renderEntry("Typ", caser.String(article.Type)) + "\n"
	strBuf += d.breakingText(article.Breaking) + "\n"
	strBuf += d.shared.style.ActiveHighlightStyle.Render("Tags:") + "\n"
	strBuf += tagStr
	if len(relatedStr) > 0 {
		strBuf += d.shared.style.ActiveHighlightStyle.Render("Verwandt:") + "\n"
		strBuf += relatedStr
	}

	return lipgloss.NewStyle().PaddingLeft(2).Render(strBuf)
}

func (d Details) renderEntry(header, content string) string {
	return d.shared.style.ActiveHighlightStyle.Render(header+": ") + d.shared.style.InactiveStyle.Render(content)
}

func (d Details) breakingText(breaking bool) string {
	if breaking {
		return d.renderEntry("Eilmeldung", "Ja")
	} else {
		return d.renderEntry("Eilmeldung", "Nein")
	}
}
