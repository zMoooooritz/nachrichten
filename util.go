package main

import (
	"os/exec"
	"runtime"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func open_url(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func ContentToText(content []Content, width int) string {
	converter := md.NewConverter("", true, nil)
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
		glamour.WithStyles(NachrichtenStyleConfig),
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
			if IsHighlighted(text, "em") {
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
			sec := IsSection(text)
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

func IsSection(text string) bool {
	return IsHighlighted(text, "strong") || IsHighlighted(text, "em")
}

func IsHighlighted(text string, tag string) bool {
	return strings.HasPrefix(text, "<"+tag+">") && strings.HasSuffix(text, "</"+tag+">")
}

func padText(text string, width int) string {
	result := ""
	split := strings.Split(text, "\n")
	for _, s := range split {
		splitLen := lipgloss.Width(s)
		result += s + strings.Repeat(" ", max(width-splitLen+1, 0)) + "\n"
	}
	return result
}

// func min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
