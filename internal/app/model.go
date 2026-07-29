package app

import (

	//External

	//Internal
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	mb "github.com/DYankee/RRipper/internal/musicbrainz"
	"github.com/DYankee/RRipper/internal/nav"
	"github.com/DYankee/RRipper/internal/screens/quantize"
	"github.com/DYankee/RRipper/internal/screens/release"
	"github.com/DYankee/RRipper/internal/screens/results"
	"github.com/DYankee/RRipper/internal/screens/search"
	"github.com/DYankee/RRipper/internal/styles"
)

type screen int

const (
	screenSearch screen = iota
	screenResults
	screenRelease
	screenQuantize
)

type Model struct {
	// global state
	current screen
	width   int
	height  int
	mb      *mb.MusicBrainz // musicbrainz client

	// screens
	search   search.Model
	results  results.Model
	release  release.Model
	quantize quantize.Model
}

func New(client *mb.MusicBrainz) Model {
	return Model{
		current: screenSearch,
		mb:      client,
		search:  search.New(client),
	}
}

func (m Model) Init() tea.Cmd {
	return m.search.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		// return m.delegateUpdate(msg)
		return m, nil

	case nav.ToSearch:
		m.current = screenSearch
		m.search = search.New(m.mb)
		return m, m.search.Init()

	case nav.ToResults:
		m.current = screenResults
		m.results = results.New(msg.Releases)
		return m, m.results.Init()

	case nav.ToAlbum:
		m.current = screenRelease
		m.release = release.New(msg.Release, m.width, m.height)
		return m, m.release.Init()

	case nav.ToQuantize:
		m.current = screenQuantize
		m.quantize = quantize.New(msg.Release, m.width, m.height)
		return m, m.quantize.Init()

	case nav.GoBack:
		return m.goBack()
	}

	return m.delegateUpdate(msg)
}

// goBack steps back one screen without reinitializing — the existing model
// instances on app.Model are still populated from when we navigated forward.
func (m Model) goBack() (tea.Model, tea.Cmd) {
	switch m.current {
	case screenResults:
		m.current = screenSearch
		m.search = search.New(m.mb) // clear the search form
	case screenRelease:
		m.current = screenResults // results.Model still intact
	case screenQuantize:
		m.current = screenRelease // album.Model still intact
	}
	return m, nil
}

func (m Model) delegateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.current {
	case screenSearch:
		m.search, cmd = m.search.Update(msg)
	case screenResults:
		m.results, cmd = m.results.Update(msg)
	case screenRelease:
		m.release, cmd = m.release.Update(msg)
	case screenQuantize:
		m.quantize, cmd = m.quantize.Update(msg)
	}
	return m, cmd
}

func (m Model) View() tea.View {
	cView := m.GetCurrentView()

	card := styles.WindowStyle.
		Render(cView.Content)

	body := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		card,
		lipgloss.WithWhitespaceStyle(
			lipgloss.NewStyle().Background(lipgloss.Black),
		),
	)
	v := tea.NewView(body)
	v.AltScreen = true
	return v
}

func (m Model) GetCurrentView() tea.View {
	switch m.current {
	case screenSearch:
		return m.search.View()
	case screenResults:
		return m.results.View()
	case screenRelease:
		return m.release.View()
	case screenQuantize:
		return m.quantize.View()
	default:
		return tea.NewView("No View Selected")
	}
}
