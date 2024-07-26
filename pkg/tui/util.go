package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
)

func NewDotSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return s
}

type LoadingNewsFailed struct{}
