package tagesschau

import "strings"

func ContentToParagraphs(content []Content) []string {
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
	return paragraphs
}

func isSection(text string) bool {
	return isHighlighted(text, "strong") || isHighlighted(text, "em")
}

func isHighlighted(text string, tag string) bool {
	return strings.HasPrefix(text, "<"+tag+">") && strings.HasSuffix(text, "</"+tag+">")
}
