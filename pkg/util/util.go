package util

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func PadText(text string, width int) string {
	result := ""
	split := strings.Split(text, "\n")
	for _, s := range split {
		splitLen := lipgloss.Width(s)
		result += s + strings.Repeat(" ", max(width-splitLen+1, 0)) + "\n"
	}
	return result
}
