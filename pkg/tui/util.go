package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
)

func NewDotSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return s
}

type LoadingNewsFailed struct{}
type LoadingArticlesFailed struct{}

type ChangedActiveArticle tagesschau.Article
type RefreshActiveViewer struct{}
type ShowTextViewer struct{}
