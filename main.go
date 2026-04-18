// main.go
package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	Internal "github.com/DYankee/RRipper/internal"
	"github.com/DYankee/RRipper/models"
)

type modelName string

const (
	searchModel  modelName = "search"
	resultsModel modelName = "results"
)

type rootModel struct {
	width        int
	height       int
	currentModel modelName
	viewMap      map[modelName]tea.Model
	mbClient     *Internal.MusicBrainz
	loading      bool
}

func newRootModel() rootModel {
	vm := make(map[modelName]tea.Model)
	vm[searchModel] = models.InitialSearchModel()
	// Results will be initialized once data is received

	return rootModel{
		currentModel: searchModel,
		viewMap:      vm,
	}
}

func (m rootModel) Init() tea.Cmd {
	return m.viewMap[m.currentModel].Init()
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "f1":
			m.currentModel = searchModel
			return m, nil
		}

	// 1. User clicked "Submit" in the Search View
	case models.SearchRequestMsg:
		m.loading = true
		return m, m.performSearch(msg.Artist, msg.Album)

		// 2. API Call Succeeded
	case models.SearchResultMsg:
		m.loading = false
		m.currentModel = resultsView
		m.viewMap[resultsView] = models.NewResultsModel(msg.Response)
		return m, nil

		// 3. API Call Failed
	case models.SearchErrorMsg:
		m.loading = false
		// You could update a status message here
		log.Printf("Error: %v", msg.Err)
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	// Pass messages to sub-models
	var cmd tea.Cmd
	if active, ok := m.viewMap[m.currentModel]; ok {
		m.viewMap[m.currentModel], cmd = active.Update(msg)
	}
	return m, cmd
}

func (m rootModel) View() tea.View {
	header := lipgloss.NewStyle().
		Width(m.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Render("RRipper v2")

	content := m.viewMap[m.currentModel].View().Content

	result := lipgloss.JoinVertical(lipgloss.Left, header, content)
	return tea.NewView(result)
}

func (m rootModel) performSearch(artist, album string) tea.Cmd {
	return func() tea.Msg {
		// You might want to make 'format' an input in your UI later,
		// for now we'll use "CD" as a placeholder.
		res, err := m.mbClient.SearchRelease(artist, album, "CD")
		if err != nil {
			return models.SearchErrorMsg{Err: err}
		}
		return models.SearchResultMsg{Response: res}
	}
}

func main() {
	p := tea.NewProgram(newRootModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
