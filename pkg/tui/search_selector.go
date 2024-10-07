package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
)

var (
	noArticles = []tagesschau.Article{
		{
			Topline: "LEER",
			Date:    time.Now(),
		},
	}
)

type SearchSelector struct {
	BaseSelector
	search textinput.Model
}

func loadArticles(seachTerm string) tea.Cmd {
	return func() tea.Msg {
		articles, err := tagesschau.SearchArticles(seachTerm)
		if err == nil {
			return articles
		}
		return LoadingArticlesFailed{}
	}
}

func NewSearchSelector(selector BaseSelector) *SearchSelector {
	searchInput := textinput.New()
	searchInput.Prompt = "> "
	searchInput.Placeholder = "Suche ..."

	selector.articles = noArticles
	return &SearchSelector{
		BaseSelector: selector,
		search:       searchInput,
	}
}

func (s SearchSelector) Init() tea.Cmd {
	return nil
}

func (s *SearchSelector) Update(msg tea.Msg) (Selector, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case LoadingArticlesFailed:
		s.articles = noArticles
		s.list.SetItems([]list.Item{})
		s.selectedIndex = 0
		cmds = append(cmds, s.PushCurrentArticle())
	case tagesschau.SearchResult:
		result := tagesschau.SearchResult(msg)
		s.articles = result.Articles
		s.rebuildList()
		if s.shared.config.Settings.PreloadThumbnails {
			go s.shared.imageCache.LoadThumbnails(s.articles)
		}
		cmds = append(cmds, s.PushCurrentArticle())
	case tea.KeyMsg:
		if s.isFocused {
			if s.shared.mode == NORMAL_MODE {
				switch {
				case key.Matches(msg, s.shared.keymap.search):
					s.shared.mode = INSERT_MODE
					s.search.Focus()
					s.search.Reset()
					return s, nil
				}
			}
			if s.shared.mode == NORMAL_MODE {
				s.list, cmd = s.list.Update(msg)
				cmds = append(cmds, cmd)
			}
			if s.shared.mode == INSERT_MODE {
				switch {
				case key.Matches(msg, s.shared.keymap.escape):
					s.shared.mode = NORMAL_MODE
					s.search.Blur()
					s.search.Reset()
				case key.Matches(msg, s.shared.keymap.confirm):
					s.shared.mode = NORMAL_MODE
					s.search.Blur()
					cmds = append(cmds, loadArticles(s.search.Value()))
				}
			}
		}
	}

	if s.isFocused {
		s.search, cmd = s.search.Update(msg)
		cmds = append(cmds, cmd)
	}

	bs, cmd := s.BaseSelector.Update(msg)
	cmds = append(cmds, cmd)
	return &SearchSelector{BaseSelector: bs, search: s.search}, tea.Batch(cmds...)
}

func (s SearchSelector) View() string {
	s.search.Width = s.width - 3

	searchView := lipgloss.NewStyle().MarginBottom(1).Render(s.search.View())

	s.list.SetSize(s.width, s.height-lipgloss.Height(searchView))

	return lipgloss.JoinVertical(lipgloss.Left, searchView, s.list.View())
}
