package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/zMoooooritz/nachrichten/pkg/config"
)

type KeyMap struct {
	quit      key.Binding
	right     key.Binding
	left      key.Binding
	up        key.Binding
	down      key.Binding
	prev      key.Binding
	next      key.Binding
	full      key.Binding
	start     key.Binding
	end       key.Binding
	image     key.Binding
	open      key.Binding
	video     key.Binding
	shortNews key.Binding
	help      key.Binding
}

func GetKeyMap(keys config.Keys) KeyMap {
	return KeyMap{
		quit:      toHelpBinding(keys.Quit, "quit"),
		right:     toHelpBinding(keys.Right, "right"),
		left:      toHelpBinding(keys.Left, "left"),
		up:        toHelpBinding(keys.Up, "up"),
		down:      toHelpBinding(keys.Down, "down"),
		next:      toHelpBinding(keys.Next, "next"),
		prev:      toHelpBinding(keys.Prev, "prev"),
		full:      toHelpBinding(keys.Full, "full"),
		start:     toHelpBinding(keys.Start, "start"),
		end:       toHelpBinding(keys.End, "end"),
		image:     toHelpBinding(keys.ToggleThumbnail, "image"),
		open:      toHelpBinding(keys.OpenArticle, "open"),
		video:     toHelpBinding(keys.OpenVideo, "video"),
		shortNews: toHelpBinding(keys.OpenShortNews, "shortnews"),
		help:      toHelpBinding(keys.Help, "help"),
	}
}

func toHelpBinding(binds []string, name string) key.Binding {
	if len(binds) == 0 {
		binds = append(binds, "NOKEY")
	}

	return key.NewBinding(
		key.WithKeys(binds...),
		key.WithHelp(keybindsToHelpText(binds), name),
	)
}

func keybindsToHelpText(keybinds []string) string {
	for i := range keybinds {
		if keybinds[i] == "up" {
			keybinds[i] = "↑"
		}
		if keybinds[i] == "down" {
			keybinds[i] = "↓"
		}
		if keybinds[i] == "left" {
			keybinds[i] = "←"
		}
		if keybinds[i] == "right" {
			keybinds[i] = "→"
		}
	}

	if len(keybinds) == 1 && keybinds[0] == "NOKEY" {
		return "NOKEY"
	}

	if len(keybinds) == 2 {
		return keybinds[0] + "/" + keybinds[1]
	}

	return keybinds[0]
}

func viewportKeymap(k config.Keys) viewport.KeyMap {
	km := viewport.DefaultKeyMap()
	km.Up = toBinding(k.Up)
	km.Down = toBinding(k.Down)
	km.PageUp = toBinding(k.Start)
	km.PageDown = toBinding(k.End)
	km.HalfPageUp = disabledBinding()
	km.HalfPageDown = disabledBinding()
	return km
}

func listKeymap(k config.Keys) list.KeyMap {
	km := list.DefaultKeyMap()
	km.CursorUp = toBinding(k.Up)
	km.CursorDown = toBinding(k.Down)
	km.GoToStart = toBinding(k.Start)
	km.GoToEnd = toBinding(k.End)
	km.PrevPage = disabledBinding()
	km.NextPage = disabledBinding()
	km.Filter = disabledBinding()
	km.ClearFilter = disabledBinding()
	km.CancelWhileFiltering = disabledBinding()
	km.AcceptWhileFiltering = disabledBinding()
	km.ShowFullHelp = disabledBinding()
	km.CloseFullHelp = disabledBinding()
	km.Quit = disabledBinding()
	km.ForceQuit = disabledBinding()
	return km
}

func toBinding(keybinds []string) key.Binding {
	return key.NewBinding(key.WithKeys(keybinds...))
}

func disabledBinding() key.Binding {
	return key.NewBinding(key.WithDisabled())
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
		{k.next},
		{k.prev},
		{k.full},
		{k.start},
		{k.end},
		{k.image},
		{k.open},
		{k.video},
		{k.shortNews},
		{k.help},
		{k.quit},
	}
}
