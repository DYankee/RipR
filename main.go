package main

import (
	"log"

	"github.com/DYankee/RRipper/views" // Update this to your actual module path

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// viewName defines the keys for our view map
type viewName string

const (
	searchView  viewName = "search"
	resultsView viewName = "results"
)

// rootModel is the "Master" model that manages sub-views and global state
type rootModel struct {
	width       int
	height      int
	currentView viewName
	viewMap     map[viewName]tea.Model
}

func newRootModel() rootModel {
	// Initialize the map with copies of your search view
	vm := make(map[viewName]tea.Model)
	vm[searchView] = views.InitialModel()
	vm[resultsView] = views.InitialModel() // Using a copy as requested

	return rootModel{
		currentView: searchView,
		viewMap:     vm,
	}
}

func (m rootModel) Init() tea.Cmd {
	// Initialize the starting view
	return m.viewMap[m.currentView].Init()
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.height = msg.Height
		m.width = msg.Width
	}

	// 1. Global Keybindings (Switching views)
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "f1": // Example: Swap to Search
			m.currentView = searchView
			return m, m.viewMap[m.currentView].Init()
		case "f2": // Example: Swap to Results
			m.currentView = resultsView
			return m, m.viewMap[m.currentView].Init()
		}
	}

	// 2. Pass messages to the ACTIVE sub-model
	activeModel := m.viewMap[m.currentView]
	newModel, newCmd := activeModel.Update(msg)

	// Update the map with the modified state of the sub-model
	m.viewMap[m.currentView] = newModel
	cmd = newCmd

	return m, cmd
}

func (m rootModel) View() tea.View {
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Render("header")
	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render("footer")

	contentView := m.viewMap[m.currentView].View()

	result := lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		contentView.Content,
		footer,
	)

	v := tea.NewView(result)
	return v
}

func main() {
	// Setup logging
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := tea.NewProgram(newRootModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
