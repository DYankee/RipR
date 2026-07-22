// main.go
package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	Internal "github.com/DYankee/RRipper/internal"
	"github.com/DYankee/RRipper/models"
	"github.com/DYankee/godacity"
)

type modelName string

const (
	searchModel  modelName = "search"
	resultsModel modelName = "results"
)

type rootModel struct {
	// size information
	fullWidth    int
	fullHeight   int
	headerHeight int
	footerHeight int
	showFooter   bool

	// Dependencies
	mbClient *Internal.MusicBrainz
	godacity *godacity.Audacity

	currentModel modelName
	viewMap      map[modelName]tea.Model
	loading      bool
}

func newRootModel() rootModel {
	vm := make(map[modelName]tea.Model)

	// Init models that don't require data
	vm[searchModel] = models.NewSearchModel()

	return rootModel{
		mbClient:     Internal.NewClient(),
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
		m.handleWindowSizeMsg(msg)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "f1":
			m.currentModel = searchModel
			return m, nil
		}

	// 1. User clicked "Submit" in the Search View
	case models.SearchSubmittedMsg:
		m.loading = true
		return m, m.performSearch(msg.Artist, msg.Album)

		// 2. API Call Succeeded
	case models.SearchResultMsg:
		m.loading = false
		m.currentModel = resultsModel
		m.viewMap[resultsModel] = models.NewResultsModel(msg.Response)
		return m, nil

		// 3. API Call Failed
	case models.SearchErrMsg:
		m.loading = false
		// You could update a status message here
		log.Printf("Error: %v", msg.Err)
		return m, nil

	}

	// Pass messages to sub-models
	if active, ok := m.viewMap[m.currentModel]; ok {
		m.viewMap[m.currentModel], cmd = active.Update(msg)
	}
	return m, cmd
}

func (m rootModel) View() tea.View {
	// Guard against the first frame where dimensions are 0
	if m.fullWidth == 0 || m.fullHeight == 0 {
		return tea.NewView("")
	}

	headerStyle := lipgloss.NewStyle().
		Width(m.fullWidth).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("240"))

	header := headerStyle.Render("RRipper v2")

	// Calculate content height: total height minus header height
	headerHeight := lipgloss.Height(header)
	contentHeight := m.fullHeight - headerHeight

	// Get the content string from the sub-model
	contentStr := m.viewMap[m.currentModel].View().Content

	// Style the content to fill the remaining height and width
	// This ensures the background color #1a1a1a fills the whole screen
	styledContent := lipgloss.NewStyle().
		Width(m.fullWidth).
		Height(contentHeight).
		Background(lipgloss.Color("#1a1a1a")).
		Render(contentStr)

	result := lipgloss.JoinVertical(lipgloss.Left, header, styledContent)

	// Create the view and ensure we pass the cursor from the sub-model
	v := tea.NewView(result)
	v.Cursor = m.viewMap[m.currentModel].View().Cursor
	return v
}

func (m *rootModel) performSearch(artist, album string) tea.Cmd {
	return func() tea.Msg {
		// You might want to make 'format' an input in your UI later,
		// for now we'll use "CD" as a placeholder.
		res, err := m.mbClient.SearchRelease(artist, album, "CD")
		if err != nil {
			return models.SearchErrMsg{Err: err}
		}
		return models.SearchResultMsg{Response: res}
	}
}

func (m *rootModel) handleWindowSizeMsg(msg tea.WindowSizeMsg) {
	m.fullHeight = msg.Height
	m.fullWidth = msg.Width
}

func main() {
	p := tea.NewProgram(newRootModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
