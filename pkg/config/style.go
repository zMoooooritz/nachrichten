package config

import (
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/lipgloss"
)

const (
	Ellipsis            = "…"
	SingleFillCharacter = "─"
	DoubleFillCharacter = "═"
)

type Style struct {
	ActiveStyle            lipgloss.Style
	TitleActiveStyle       lipgloss.Style
	ListHeaderActiveStyle  lipgloss.Style
	ListActiveStyle        lipgloss.Style
	ReaderTitleActiveStyle lipgloss.Style
	ReaderInfoActiveStyle  lipgloss.Style

	InactiveStyle            lipgloss.Style
	TitleInactiveStyle       lipgloss.Style
	ListHeaderInactiveStyle  lipgloss.Style
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

	listBaseStyle := lipgloss.NewStyle().Padding(0, 1, 1, 1).Margin(1, 1, 0, 1).BorderLeft(true).BorderBottom(true).BorderRight(true).BorderTop(false)

	s.ActiveStyle = lipgloss.NewStyle().Foreground(primaryColor).BorderForeground(primaryColor)
	s.TitleActiveStyle = lipgloss.NewStyle().Background(primaryColor).Foreground(secondaryColor)
	s.ListActiveStyle = listBaseStyle.Copy().BorderForeground(primaryColor).BorderStyle(lipgloss.DoubleBorder())
	s.ListHeaderActiveStyle = func() lipgloss.Style {
		b := lipgloss.DoubleBorder()
		b.Right = "╠"
		b.Left = "╣"
		return s.ActiveStyle.Copy().BorderStyle(b).Padding(0, 1)
	}()
	s.ReaderTitleActiveStyle = func() lipgloss.Style {
		b := lipgloss.DoubleBorder()
		b.Right = "╠"
		return s.ActiveStyle.Copy().BorderStyle(b).Padding(0, 1)
	}()
	s.ReaderInfoActiveStyle = func() lipgloss.Style {
		b := lipgloss.DoubleBorder()
		b.Left = "╣"
		return s.ReaderTitleActiveStyle.Copy().BorderStyle(b)
	}()

	s.InactiveStyle = lipgloss.NewStyle().Foreground(secondaryColor).BorderForeground(secondaryColor)
	s.TitleInactiveStyle = lipgloss.NewStyle().Foreground(secondaryColor)
	s.ListInactiveStyle = listBaseStyle.Copy().BorderForeground(secondaryColor).BorderStyle(lipgloss.RoundedBorder())
	s.ListHeaderInactiveStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		b.Left = "┤"
		return s.InactiveStyle.Copy().BorderStyle(b).Padding(0, 1)
	}()
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
