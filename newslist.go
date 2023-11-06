package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/truncate"
)

type NewsDelegate struct {
	Styles  Style
	height  int
	spacing int
}

func NewNewsDelegate(s Style) NewsDelegate {
	return NewsDelegate{
		Styles:  s,
		height:  2,
		spacing: 1,
	}
}

func (n NewsDelegate) Height() int {
	return n.height
}

func (n NewsDelegate) Spacing() int {
	return n.spacing
}

func (n NewsDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (n NewsDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		title, desc string
		s           = &n.Styles
	)

	entry, ok := item.(NewsEntry)
	if ok {
		title = entry.Title()
		desc = entry.Description()
	} else {
		return
	}

	if m.Width() <= 0 {
		// short-circuit
		return
	}

	// Prevent text from exceeding list width
	textwidth := uint(m.Width() - s.ItemNormalTitle.GetPaddingLeft() - s.ItemNormalTitle.GetPaddingRight())
	title = truncate.StringWithTail(title, textwidth, ellipsis)
	var lines []string
	for i, line := range strings.Split(desc, "\n") {
		if i >= n.height-1 {
			break
		}
		lines = append(lines, truncate.StringWithTail(line, textwidth, ellipsis))
	}
	desc = strings.Join(lines, "\n")

	var isSelected = index == m.Index()

	if entry.Breaking && !isSelected {
		title = s.ItemBreakingTitle.Render(title)
		desc = s.ItemBreakingDesc.Render(desc)
	}

	if isSelected {
		title = s.ItemSelectedTitle.Render(title)
		desc = s.ItemSelectedDesc.Render(desc)
	} else {
		title = s.ItemNormalTitle.Render(title)
		desc = s.ItemNormalDesc.Render(desc)
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}
