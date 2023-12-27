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
	TextHighlightStyle lipgloss.Style

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

	ActiveTabBorder   lipgloss.Border
	InactiveTabBorder lipgloss.Border

	ReaderStyle ansi.StyleConfig
}

func GruvboxTheme() Theme {
	return Theme{
		PrimaryColor:         "#EBDBB2",
		ShadedColor:          "#928374",
		HighlightColor:       "#458588",
		HighlightShadedColor: "#83A598",
		WarningColor:         "#FB4934",
		WarningShadedColor:   "#CC241D",
		ReaderHighlightColor: "#FABD2F",
		ReaderHeadingColor:   "#8EC07C",
	}
}

func NewsStyle(t Theme) (s Style) {
	primaryColor := lipgloss.Color(t.PrimaryColor)
	shadedColor := lipgloss.Color(t.ShadedColor)
	highlightColor := lipgloss.Color(t.HighlightColor)
	highlightShadedColor := lipgloss.Color(t.HighlightShadedColor)
	warningColor := lipgloss.Color(t.WarningColor)
	warningShadedColor := lipgloss.Color(t.WarningShadedColor)

	listBaseStyle := lipgloss.NewStyle().Padding(0, 1, 1, 1).Margin(0, 1, 0, 1)

	s.TextHighlightStyle = lipgloss.NewStyle().Foreground(highlightColor).BorderForeground(primaryColor)

	s.ActiveStyle = lipgloss.NewStyle().Foreground(highlightColor).BorderForeground(highlightColor)
	s.TitleActiveStyle = lipgloss.NewStyle().Background(highlightColor).Foreground(primaryColor)
	s.ListActiveStyle = listBaseStyle.Copy().BorderForeground(highlightColor).BorderStyle(lipgloss.DoubleBorder())
	s.ListHeaderActiveStyle = s.ActiveStyle.Copy().Padding(0, 1).Border(lipgloss.DoubleBorder(), false, false, true, false)
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

	s.InactiveStyle = lipgloss.NewStyle().Foreground(primaryColor).BorderForeground(primaryColor)
	s.TitleInactiveStyle = lipgloss.NewStyle().Foreground(primaryColor)
	s.ListInactiveStyle = listBaseStyle.Copy().BorderForeground(primaryColor).BorderStyle(lipgloss.RoundedBorder())
	s.ListHeaderInactiveStyle = s.InactiveStyle.Copy().Padding(0, 1).Border(lipgloss.RoundedBorder(), false, false, true, false)
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

	s.ItemNormalTitle = lipgloss.NewStyle().Foreground(primaryColor).Padding(0, 0, 0, 2)

	s.ItemNormalDesc = s.ItemNormalTitle.Copy().Foreground(shadedColor)

	s.ItemSelectedTitle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(highlightColor).Foreground(highlightShadedColor).Padding(0, 0, 0, 1)

	s.ItemSelectedDesc = s.ItemSelectedTitle.Copy().Foreground(highlightColor)

	s.ItemBreakingTitle = lipgloss.NewStyle().Foreground(warningColor).Padding(0, 0, 0, 2)

	s.ItemBreakingDesc = s.ItemBreakingTitle.Copy().Foreground(warningShadedColor)

	s.ActiveTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	s.InactiveTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	s.ReaderStyle = CreateReaderStyle(t)

	return s
}
