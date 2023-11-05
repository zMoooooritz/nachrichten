package main

import "github.com/charmbracelet/lipgloss"

const (
	ellipsis = "…"
)

type Style struct {
	ActiveStyle            lipgloss.Style
	TitleActiveStyle       lipgloss.Style
	ListActiveStyle        lipgloss.Style
	ReaderTitleActiveStyle lipgloss.Style
	ReaderInfoActiveStyle  lipgloss.Style

	InactiveStyle            lipgloss.Style
	TitleInactiveStyle       lipgloss.Style
	ListInactiveStyle        lipgloss.Style
	ReaderTitleInactiveStyle lipgloss.Style
	ReaderInfoInactiveStyle  lipgloss.Style

	// the normal item state
	ItemNormalTitle lipgloss.Style
	ItemNormalDesc  lipgloss.Style

	// he selected item state
	ItemSelectedTitle lipgloss.Style
	ItemSelectedDesc  lipgloss.Style

	// The breaking item state
	ItemBreakingTitle lipgloss.Style
	ItemBreakingDesc  lipgloss.Style
}

func DefaultNewsStyle() (s Style) {
	primaryColor := lipgloss.Color("62")
	secondaryColor := lipgloss.Color("230")
	normalTitleColor := lipgloss.Color("#DDDDDD")
	normalDescColor := lipgloss.Color("#777777")
	selectedPrimaryColor := lipgloss.Color("#AD58B4")
	selectedSecondaryColor := lipgloss.Color("#EE6FF8")
	breakingColor := lipgloss.Color("#FF0000")

	s.ActiveStyle = lipgloss.NewStyle().Foreground(primaryColor).BorderForeground(primaryColor)
	s.TitleActiveStyle = lipgloss.NewStyle().Background(primaryColor).Foreground(secondaryColor)
	s.ListActiveStyle = s.ActiveStyle.Copy().Padding(1, 1, 1, 1).Margin(0, 1, 0, 1).BorderStyle(lipgloss.RoundedBorder())
	s.ReaderTitleActiveStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return s.ActiveStyle.Copy().BorderStyle(b).Padding(0, 1)
	}()
	s.ReaderInfoActiveStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return s.ReaderTitleActiveStyle.Copy().BorderStyle(b)
	}()

	s.InactiveStyle = lipgloss.NewStyle()
	s.TitleInactiveStyle = lipgloss.NewStyle().Foreground(secondaryColor)
	s.ListInactiveStyle = s.InactiveStyle.Copy().Padding(1, 1, 1, 1).Margin(0, 1, 0, 1).BorderStyle(lipgloss.RoundedBorder())
	s.ReaderTitleInactiveStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return s.InactiveStyle.Copy().BorderStyle(b).Padding(0, 1)
	}()
	s.ReaderInfoInactiveStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return s.ReaderTitleInactiveStyle.Copy().BorderStyle(b)
	}()

	s.ItemNormalTitle = lipgloss.NewStyle().Foreground(normalTitleColor).Padding(0, 0, 0, 2)

	s.ItemNormalDesc = s.ItemNormalTitle.Copy().Foreground(normalDescColor)

	s.ItemSelectedTitle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(selectedPrimaryColor).Foreground(selectedSecondaryColor).Padding(0, 0, 0, 1)

	s.ItemSelectedDesc = s.ItemSelectedTitle.Copy().Foreground(selectedPrimaryColor)

	s.ItemBreakingTitle = lipgloss.NewStyle().Foreground(breakingColor).Padding(0, 0, 0, 2)

	s.ItemBreakingDesc = s.ItemBreakingTitle.Copy().Foreground(breakingColor)

	return s
}
