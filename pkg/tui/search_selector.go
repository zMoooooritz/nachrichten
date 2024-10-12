package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
)

type SearchSelector struct {
	BaseSelector
	search textinput.Model
}

func loadArticles(seachTerm string) tea.Cmd {
	return func() tea.Msg {
		searchResult, err := tagesschau.SearchArticles(seachTerm)
		if err == nil && len(searchResult.Articles) > 0 {
			return searchResult
		}
		return LoadingArticlesFailed{}
	}
}

func NewSearchSelector(selector BaseSelector) *SearchSelector {
	searchInput := textinput.New()
	searchInput.Prompt = ""
	searchInput.Placeholder = "Suche ..."
	searchInput.PromptStyle = selector.shared.style.ItemSelectedTitle
	searchInput.Cursor.Style = selector.shared.style.InactiveStyle
	searchInput.Cursor.TextStyle = selector.shared.style.InactiveStyle
	searchInput.TextStyle = selector.shared.style.InactiveStyle

	selector.articles = []tagesschau.Article{tagesschau.EMPTY_ARTICLE()}
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
		s.articles = []tagesschau.Article{tagesschau.EMPTY_ARTICLE()}
		s.list.SetItems([]list.Item{})
		s.selectedIndex = 0
		cmds = append(cmds, s.PushSelectedArticle())
	case tagesschau.SearchResult:
		result := tagesschau.SearchResult(msg)
		s.articles = result.Articles
		s.rebuildList()
		if s.shared.config.Settings.PreloadThumbnails {
			go s.shared.imageCache.LoadThumbnails(s.articles)
		}
		cmds = append(cmds, s.PushSelectedArticle())
	case tea.KeyMsg:
		if s.isFocused && s.isVisible {
			if s.shared.mode == NORMAL_MODE {
				switch {
				case key.Matches(msg, s.shared.keymap.search):
					s.shared.mode = INSERT_MODE
					s.search.Focus()
					s.search.Reset()
					bs, cmd := s.BaseSelector.Update(msg)
					cmds = append(cmds, cmd)
					return &SearchSelector{BaseSelector: bs, search: s.search}, tea.Batch(cmds...)
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

	if s.isFocused && s.isVisible {
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
