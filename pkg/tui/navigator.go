package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	headerText         string = "Nachrichten"
	regionalHeaderText string = "Regional"
	nationalHeaderText string = "National"
	searchHeaderText   string = "Suche"

	selectorCount int = 3
)

type Navigator struct {
	shared              *SharedState
	selectors           []Selector
	activeSelectorIndex int
	isFocused           bool
	isVisible           bool
	width               int
	height              int
}

func NewNavigator(shared *SharedState) *Navigator {
	selectors := []Selector{
		NewHomeSelector(NewSelector(ST_NATIONAL, shared, true)),
		NewHomeSelector(NewSelector(ST_REGIONAL, shared, false)),
		NewSearchSelector(NewSelector(ST_SEARCH, shared, false)),
	}

	return &Navigator{
		shared:              shared,
		selectors:           selectors,
		activeSelectorIndex: 0,
		isFocused:           true,
		isVisible:           true,
	}
}

func (n *Navigator) nextSelector() {
	n.gotoSelector((n.activeSelectorIndex + 1) % selectorCount)
}

func (n *Navigator) prevSelector() {
	n.gotoSelector((selectorCount + n.activeSelectorIndex - 1) % selectorCount)
}

func (n *Navigator) selectSearchSelector() {
	n.gotoSelector(selectorCount - 1)
}

func (n *Navigator) gotoSelector(index int) {
	n.selectors[n.activeSelectorIndex].SetVisible(false)
	n.selectors[n.activeSelectorIndex].SetActive(false)
	n.selectors[n.activeSelectorIndex].SetFocused(false)
	n.activeSelectorIndex = index
	n.selectors[n.activeSelectorIndex].SetVisible(true)
	n.selectors[n.activeSelectorIndex].SetActive(true)
	n.selectors[n.activeSelectorIndex].SetFocused(true)
}

func (n *Navigator) SetDims(w, h int) {
	n.width = w
	n.height = h
}

func (n Navigator) Init() tea.Cmd {
	return nil
}

func (n *Navigator) Update(msg tea.Msg) (*Navigator, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if n.shared.mode == INSERT_MODE {
			break
		}

		switch {
		case key.Matches(msg, n.shared.keymap.next):
			if n.isFocused && n.isVisible {
				n.nextSelector()
				cmds = append(cmds, n.selectors[n.activeSelectorIndex].PushSelectedArticle())
			}
		case key.Matches(msg, n.shared.keymap.prev):
			if n.isFocused && n.isVisible {
				n.prevSelector()
				cmds = append(cmds, n.selectors[n.activeSelectorIndex].PushSelectedArticle())
			}
		case key.Matches(msg, n.shared.keymap.right):
			if n.isVisible {
				n.isFocused = false
			}
		case key.Matches(msg, n.shared.keymap.left):
			if n.isVisible {
				n.isFocused = true
			}
		case key.Matches(msg, n.shared.keymap.full):
			n.isVisible = !n.isVisible
		case key.Matches(msg, n.shared.keymap.search):
			if n.isFocused && n.isVisible {
				n.selectSearchSelector()
			}
		}
	}

	var updatedSelector Selector
	for i, selector := range n.selectors {
		updatedSelector, cmd = selector.Update(msg)
		n.selectors[i] = updatedSelector
		cmds = append(cmds, cmd)
	}

	return n, tea.Batch(cmds...)
}

func (n Navigator) View() string {
	if !n.isVisible {
		return ""
	}

	headerView := n.headerView()
	tabView := n.tabView([]string{nationalHeaderText, regionalHeaderText, searchHeaderText}, n.activeSelectorIndex)

	style := n.shared.style.ListInactiveStyle
	if n.isFocused {
		style = n.shared.style.ListActiveStyle
	}

	n.selectors[n.activeSelectorIndex].SetDims(n.width, n.height-lipgloss.Height(headerView)-lipgloss.Height(tabView)-lipgloss.Height(style.Render(""))+1)

	return lipgloss.JoinVertical(lipgloss.Left, headerView, style.Render(lipgloss.JoinVertical(lipgloss.Left, tabView, n.selectors[n.activeSelectorIndex].View())))
}

func (n Navigator) headerView() string {
	headerStyle := n.shared.style.ListHeaderInactiveStyle
	if n.isFocused {
		headerStyle = n.shared.style.ListHeaderActiveStyle
	}

	centeredHeader := lipgloss.PlaceHorizontal(n.width, lipgloss.Center, headerText)
	return headerStyle.Render(centeredHeader)
}

func (n Navigator) tabView(names []string, activeIndex int) string {
	cellWidth := (n.width - 2*len(names)) / len(names)
	var widths []int
	for i := 0; i < len(names)-1; i++ {
		widths = append(widths, cellWidth)
	}
	widths = append(widths, n.width-(len(names)-1)*cellWidth-2*len(names))

	result := ""
	for i, name := range names {
		border := n.shared.style.InactiveTabBorder
		style := n.shared.style.InactiveStyle
		if i == activeIndex {
			border = n.shared.style.ActiveTabBorder
			style = n.shared.style.TextHighlightStyle
		}
		centeredText := lipgloss.PlaceHorizontal(widths[i], lipgloss.Center, name)
		result = lipgloss.JoinHorizontal(lipgloss.Center, result, style.MarginBottom(1).BorderStyle(border).Render(centeredText))
	}
	return result
}
