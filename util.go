package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

func openUrl(t ResourceType, c Configuration, url string) error {
	var appConfig ApplicationConfig

	switch t {
	case TypeImage:
		appConfig = c.AppConfig.Image
	case TypeAudio:
		appConfig = c.AppConfig.Audio
	case TypeVideo:
		appConfig = c.AppConfig.Video
	case TypeHTML:
		appConfig = c.AppConfig.HTML
	default:
		return defaultOpenUrl(url)
	}

	cConfig := appConfig
	cConfig.Args = append([]string(nil), appConfig.Args...)

	if cConfig.Path == "" || len(cConfig.Args) == 0 {
		return defaultOpenUrl(url)
	}

	for i, arg := range cConfig.Args {
		if arg == "$" {
			cConfig.Args[i] = url
		}
	}
	return exec.Command(cConfig.Path, cConfig.Args...).Start()
}

func defaultOpenUrl(url string) error {
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

func loadConfig(path string) Configuration {
	var config Configuration

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return config
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
	}
	return config
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
