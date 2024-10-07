package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewManager struct {
	shared            *SharedState
	viewers           []Viewer
	activeViewerIndex int
	width             int
	height            int
}

func NewViewManager(shared *SharedState) *ViewManager {
	viewers := []Viewer{
		NewReader(NewViewer(VT_TEXT, shared, true)),
		NewImageViewer(NewViewer(VT_IMAGE, shared, false)),
		NewDetails(NewViewer(VT_DETAILS, shared, false)),
	}

	return &ViewManager{
		shared:            shared,
		viewers:           viewers,
		activeViewerIndex: 0,
	}
}

func (v *ViewManager) SetDims(w, h, splitOffset int) {
	v.width = w
	v.height = h

	isViewerFullscreen := v.activeViewer().IsFullScreen()
	for _, viewer := range v.viewers {
		if isViewerFullscreen {
			viewer.SetDims(v.width, v.height)
		} else {
			viewer.SetDims(v.width-splitOffset, v.height)
		}
	}
}

func (v ViewManager) activeViewer() Viewer {
	for _, viewer := range v.viewers {
		if viewer.IsActive() {
			return viewer
		}
	}
	return v.viewers[0]
}

func (v *ViewManager) showViewer(vt ViewerType) tea.Cmd {
	currViewer := v.activeViewer()
	nextViewer := v.activeViewer()
	for _, viewer := range v.viewers {
		if viewer.ViewerType() == vt {
			nextViewer = viewer
		}
	}
	if currViewer.ViewerType() == nextViewer.ViewerType() {
		return nil
	}

	nextViewer.SetActive(true)
	nextViewer.SetFocused(currViewer.IsFocused())
	nextViewer.SetFullScreen(currViewer.IsFullScreen())

	currViewer.SetActive(false)
	currViewer.SetFocused(false)
	currViewer.SetFullScreen(false)

	return refreshFunc
}

func (v ViewManager) Init() tea.Cmd {
	return nil
}

func (v *ViewManager) Update(msg tea.Msg) (*ViewManager, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case ShowTextViewer:
		v.showViewer(VT_TEXT)
	case tea.KeyMsg:
		if v.shared.mode == NORMAL_MODE {
			switch {
			case key.Matches(msg, v.shared.keymap.article):
				v.showViewer(VT_TEXT)
			case key.Matches(msg, v.shared.keymap.image):
				v.showViewer(VT_IMAGE)
			case key.Matches(msg, v.shared.keymap.details):
				v.showViewer(VT_DETAILS)
			}
		}
	}

	var updatedViewer Viewer
	for i, viewer := range v.viewers {
		updatedViewer, cmd = viewer.Update(msg)
		v.viewers[i] = updatedViewer
		cmds = append(cmds, cmd)
	}

	return v, tea.Batch(cmds...)
}

func (v ViewManager) View() string {
	return v.viewers[v.activeViewerIndex].View()
}
