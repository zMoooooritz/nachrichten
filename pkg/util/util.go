package util

import (
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/Nachrichten/pkg/config"
	"github.com/zMoooooritz/Nachrichten/pkg/tagesschau"
)

func ContentToText(content []tagesschau.Content, width int) string {
	converter := md.NewConverter("", true, nil)
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
		glamour.WithStyles(config.NachrichtenStyleConfig),
	)

	prevType := "text"
	prevSection := false
	paragraph := ""
	var paragraphs []string
	for _, c := range content {
		switch c.Type {
		case "text":
			fallthrough
		case "headline":
			text := c.Value
			if strings.Trim(text, " ") == "" {
				continue
			}

			// drop author information
			if isHighlighted(text, "em") {
				continue
			}

			// remove hyperlinks from text
			for {
				startIdx := strings.Index(text, "<a ")
				if startIdx == -1 {
					break
				}
				endIndex := strings.Index(text, "\">")
				if endIndex == -1 {
					break
				}
				text = text[:startIdx+2] + text[endIndex+1:]
			}
			sec := isSection(text)
			if (prevType != c.Type || sec || prevSection) && paragraph != "" {
				paragraphs = append(paragraphs, paragraph)
				paragraph = ""
			}
			paragraph += text + " "
			prevSection = sec
		}
		prevType = c.Type
	}
	paragraphs = append(paragraphs, paragraph)

	result := ""
	for _, p := range paragraphs {
		text, _ := converter.ConvertString(p)
		text, _ = renderer.Render(text)
		result += text
	}
	return padText(result, width)
}

func isSection(text string) bool {
	return isHighlighted(text, "strong") || isHighlighted(text, "em")
}

func isHighlighted(text string, tag string) bool {
	return strings.HasPrefix(text, "<"+tag+">") && strings.HasSuffix(text, "</"+tag+">")
}

func padText(text string, width int) string {
	result := ""
	split := strings.Split(text, "\n")
	for _, s := range split {
		splitLen := lipgloss.Width(s)
		result += s + strings.Repeat(" ", Max(width-splitLen+1, 0)) + "\n"
	}
	return result
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
