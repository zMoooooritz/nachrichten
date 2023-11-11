package config

import (
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/lipgloss"
)

const (
	Ellipsis = "…"
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

	ReaderStyle ansi.StyleConfig
}

func DefaultThemeConfiguration() ThemeConfig {
	return ThemeConfig{
		PrimaryColor:           "#3636D1",
		SecondaryColor:         "#F8F8F2",
		NormalTitleColor:       "#DDDDDD",
		NormalDescColor:        "#777777",
		SelectedPrimaryColor:   "#AD58B4",
		SelectedSecondaryColor: "#EE6FF8",
		BreakingColor:          "#FF0000",
		ReaderHighlightColor:   "#FFB86C",
		ReaderHeadingColor:     "#BD93F9",
	}
}

func NewsStyle(t ThemeConfig) (s Style) {
	primaryColor := lipgloss.Color(t.PrimaryColor)
	secondaryColor := lipgloss.Color(t.SecondaryColor)
	normalTitleColor := lipgloss.Color(t.NormalTitleColor)
	normalDescColor := lipgloss.Color(t.NormalDescColor)
	selectedPrimaryColor := lipgloss.Color(t.SelectedPrimaryColor)
	selectedSecondaryColor := lipgloss.Color(t.SelectedSecondaryColor)
	breakingColor := lipgloss.Color(t.BreakingColor)

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

	s.ReaderStyle = CreateReaderStyle(t)

	return s
}
