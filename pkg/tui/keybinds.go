package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	quit      key.Binding
	right     key.Binding
	left      key.Binding
	next      key.Binding
	up        key.Binding
	down      key.Binding
	prev      key.Binding
	start     key.Binding
	end       key.Binding
	open      key.Binding
	video     key.Binding
	shortNews key.Binding
	help      key.Binding
}

func GetKeyMap() KeyMap {
	return KeyMap{
		quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("→/l", "right"),
		),
		left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("←/h", "left"),
		),
		up: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "down"),
		),
		down: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "up"),
		),
		start: key.NewBinding(
			key.WithKeys("g", "home"),
			key.WithHelp("g/home", "start"),
		),
		end: key.NewBinding(
			key.WithKeys("G", "end"),
			key.WithHelp("G/end", "end"),
		),
		next: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
		),
		prev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev"),
		),
		open: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open"),
		),
		video: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "video"),
		),
		shortNews: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "shortnews"),
		),
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.left, k.right, k.up, k.down, k.next, k.help, k.quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.left},
		{k.right},
		{k.up},
		{k.down},
		{k.start},
		{k.end},
		{k.next},
		{k.prev},
		{k.open},
		{k.video},
		{k.shortNews},
		{k.help},
		{k.quit},
	}
}
