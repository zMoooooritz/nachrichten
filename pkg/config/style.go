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
	TextHighlightStyle  lipgloss.Style
	ScreenCenteredStyle func(width, height int) lipgloss.Style

	ActiveStyle            lipgloss.Style
	ActiveHighlightStyle   lipgloss.Style
	HighlightStyle         lipgloss.Style
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

func topDoubleBorder() lipgloss.Border {
	border := lipgloss.DoubleBorder()
	border.BottomRight = border.MiddleRight
	border.BottomLeft = border.MiddleLeft
	return border
}

func topRoundedBorder() lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomRight = border.MiddleRight
	border.BottomLeft = border.MiddleLeft
	return border
}

func NewsStyle(t Theme) (s Style) {
	primaryColor := lipgloss.Color(t.PrimaryColor)
	shadedColor := lipgloss.Color(t.ShadedColor)
	highlightColor := lipgloss.Color(t.HighlightColor)
	highlightShadedColor := lipgloss.Color(t.HighlightShadedColor)
	warningColor := lipgloss.Color(t.WarningColor)
	warningShadedColor := lipgloss.Color(t.WarningShadedColor)
	markerColor := lipgloss.Color(t.ReaderHighlightColor)

	listBaseStyle := lipgloss.NewStyle().Padding(0, 1, 1, 1).Margin(0, 1, 0, 1)

	s.TextHighlightStyle = lipgloss.NewStyle().Foreground(highlightColor).BorderForeground(primaryColor)
	s.ScreenCenteredStyle = func(width, height int) lipgloss.Style {
		return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Height(height).AlignVertical(lipgloss.Center)
	}

	s.ActiveStyle = lipgloss.NewStyle().Foreground(highlightColor).BorderForeground(highlightColor)
	s.ActiveHighlightStyle = lipgloss.NewStyle().Foreground(highlightShadedColor).BorderForeground(highlightShadedColor).Bold(true)
	s.HighlightStyle = lipgloss.NewStyle().Foreground(markerColor).BorderForeground(markerColor).Bold(true)
	s.TitleActiveStyle = lipgloss.NewStyle().Background(highlightColor).Foreground(primaryColor)
	s.ListActiveStyle = listBaseStyle.BorderForeground(highlightColor).Border(lipgloss.DoubleBorder(), false, true, true, true)
	s.ListHeaderActiveStyle = s.ActiveStyle.Padding(0, 1).Margin(0, 1).BorderStyle(topDoubleBorder())
	s.ReaderTitleActiveStyle = func() lipgloss.Style {
		b := lipgloss.DoubleBorder()
		b.Right = "╠"
		return s.ActiveStyle.BorderStyle(b).Padding(0, 1)
	}()
	s.ReaderInfoActiveStyle = func() lipgloss.Style {
		b := lipgloss.DoubleBorder()
		b.Left = "╣"
		return s.ReaderTitleActiveStyle.BorderStyle(b)
	}()

	s.InactiveStyle = lipgloss.NewStyle().Foreground(primaryColor).BorderForeground(primaryColor)
	s.TitleInactiveStyle = lipgloss.NewStyle().Foreground(primaryColor)
	s.ListInactiveStyle = listBaseStyle.BorderForeground(primaryColor).Border(lipgloss.RoundedBorder(), false, true, true, true)
	s.ListHeaderInactiveStyle = s.InactiveStyle.Padding(0, 1).Margin(0, 1).BorderStyle(topRoundedBorder())
	s.ReaderTitleInactiveStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return s.InactiveStyle.BorderStyle(b).Padding(0, 1)
	}()
	s.ReaderInfoInactiveStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return s.ReaderTitleInactiveStyle.BorderStyle(b)
	}()

	s.ItemNormalTitle = lipgloss.NewStyle().Foreground(primaryColor).Padding(0, 0, 0, 2)

	s.ItemNormalDesc = s.ItemNormalTitle.Foreground(shadedColor)

	s.ItemSelectedTitle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(highlightColor).Foreground(highlightShadedColor).Padding(0, 0, 0, 1)

	s.ItemSelectedDesc = s.ItemSelectedTitle.Foreground(highlightColor)

	s.ItemBreakingTitle = lipgloss.NewStyle().Foreground(warningColor).Padding(0, 0, 0, 2)

	s.ItemBreakingDesc = s.ItemBreakingTitle.Foreground(warningShadedColor)

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
