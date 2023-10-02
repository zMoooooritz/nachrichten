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
	s.ActiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).BorderForeground(lipgloss.Color("62"))
	s.TitleActiveStyle = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230"))
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
	s.TitleInactiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
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

	s.ItemNormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2)

	s.ItemNormalDesc = s.ItemNormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.ItemSelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
		Padding(0, 0, 0, 1)

	s.ItemSelectedDesc = s.ItemSelectedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})

	s.ItemBreakingTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}).
		Padding(0, 0, 0, 2)

	s.ItemBreakingDesc = s.ItemBreakingTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"})

	return s
}
