package tagesschau

import (
	"regexp"
	"strings"
	"unicode"
)

const (
	timeRegex = `\b\d{1,2}:\d{2}\b`
)

func ContentToParagraphs(content []Content) []string {
	prevType := "text"
	prevSection := false
	paragraph := ""
	var paragraphs []string
	for i, c := range content {
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
			if (prevType != c.Type || sec || prevSection) && paragraph != "" || i == len(content)-1 {
				paragraphs = append(paragraphs, clean(paragraph))
				paragraph = ""
			}
			paragraph += text + " "
			prevSection = sec
		}
		prevType = c.Type
	}
	paragraphs = append(paragraphs, clean(formatLastLine(paragraph)))
	return paragraphs
}

func clean(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, s)
}

func isSection(text string) bool {
	return isHighlighted(text, "strong") || isHighlighted(text, "em")
}

func isHighlighted(text string, tag string) bool {
	return strings.HasPrefix(text, "<"+tag+">") && strings.HasSuffix(text, "</"+tag+">")
}

func formatLastLine(text string) string {
	text = strings.TrimSpace(text)

	if !(strings.HasPrefix(text, "<strong>") || containsTime(text)) || isHighlighted(text, "strong") {
		return text
	}

	text = strings.ReplaceAll(text, "<strong>", "")
	text = strings.ReplaceAll(text, "</strong>", "")
	text = strings.ReplaceAll(text, "<br />", "")
	start, end, _ := strings.Cut(text, ":")
	start = "<strong>" + start + ":" + "</strong>"
	end = strings.TrimSpace(end)

	return start + " " + end
}

func containsTime(text string) bool {
	re := regexp.MustCompile(timeRegex)
	return re.Match([]byte(text))
}
