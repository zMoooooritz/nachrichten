package tui

import (
	"strconv"

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
	pageUp    key.Binding
	pageDown  key.Binding
	search    key.Binding
	confirm   key.Binding
	escape    key.Binding
	article   key.Binding
	image     key.Binding
	details   key.Binding
	open      key.Binding
	video     key.Binding
	shortNews key.Binding
	help      key.Binding
	number    []key.Binding
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
		pageUp:    toHelpBinding(keys.PageUp, "pageup"),
		pageDown:  toHelpBinding(keys.PageDown, "pagedown"),
		search:    toHelpBinding(keys.Search, "search"),
		confirm:   toHelpBinding(keys.Confirm, "confirm"),
		escape:    toHelpBinding(keys.Escape, "escape"),
		article:   toHelpBinding(keys.ShowArticle, "article"),
		image:     toHelpBinding(keys.ShowThumbnail, "image"),
		details:   toHelpBinding(keys.ShowDetails, "details"),
		open:      toHelpBinding(keys.OpenArticle, "open"),
		video:     toHelpBinding(keys.OpenVideo, "video"),
		shortNews: toHelpBinding(keys.OpenShortNews, "shortnews"),
		help:      toHelpBinding(keys.Help, "help"),
		number:    getNumberBinds(),
	}
}

func getNumberBinds() []key.Binding {
	binds := []key.Binding{}
	for i := range 10 {
		binds = append(binds, key.NewBinding(key.WithKeys(strconv.Itoa(i))))
	}
	return binds
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

func keybindsToHelpText(binds []string) string {
	keybinds := []string{}

	for i := range binds {
		keybinds = append(keybinds, binds[i])

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

func ViewportKeymap(k config.Keys) viewport.KeyMap {
	km := viewport.DefaultKeyMap()
	km.Up = toBinding(k.Up)
	km.Down = toBinding(k.Down)
	km.PageUp = toBinding(k.PageUp)
	km.PageDown = toBinding(k.PageDown)
	km.HalfPageUp = disabledBinding()
	km.HalfPageDown = disabledBinding()
	return km
}

func ListKeymap(k config.Keys) list.KeyMap {
	km := list.DefaultKeyMap()
	km.CursorUp = toBinding(k.Up)
	km.CursorDown = toBinding(k.Down)
	km.GoToStart = toBinding(k.Start)
	km.GoToEnd = toBinding(k.End)
	km.PrevPage = toBinding(k.PageUp)
	km.NextPage = toBinding(k.PageDown)
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
		{k.left, k.right, k.up, k.down, k.prev, k.next, k.help, k.quit},
		{k.full, k.start, k.end, k.article, k.image, k.details, k.open, k.video, k.shortNews},
	}
}
