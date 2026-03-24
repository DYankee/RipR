package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michiwend/gomusicbrainz"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case errMsg:
		log.Printf("error: %v", msg.err)
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case searchRes:
		m.buildReleaseTable(msg)
		m.searchRes = msg
		m.focusIndex = 0
		m.currentView = viewSearchResult
		return m, nil

	case gomusicbrainz.Release:
		m.releaseData = msg
		m.buildReleaseResTable()
		m.currentView = viewReleaseResult
		return m, nil
	}

	return m, cmd
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.currentView {
	case viewSearch:
		return m.handleSearchInput(msg)
	case viewSearchResult:
		return m.handleSearchResultInput(msg)
	case viewReleaseResult:
		return m.handleReleaseResultInput(msg)
	}

	return m, nil
}

func (m model) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "shift+tab", "enter", "up", "down":
		if msg.String() == "enter" && m.focusIndex == len(m.inputs) {
			log.Printf(
				"Artist: %s Release: %s",
				m.inputs[0].Value(),
				m.inputs[1].Value(),
			)
			m.currentView = viewLoading
			return m, m.searchRelease()
		}

		if msg.String() == "up" || msg.String() == "shift+tab" {
			m.focusIndex--
		} else {
			m.focusIndex++
		}

		if m.focusIndex > len(m.inputs) {
			m.focusIndex = 0
		} else if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs)
		}
	}

	return m, m.updateInputs(msg)
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	for i := range m.inputs {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
		}
	}

	return tea.Batch(cmds...)
}

func (m model) handleSearchResultInput(
	msg tea.KeyMsg,
) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		return m, m.getReleaseData()
	case "esc":
		m.currentView = viewSearch
	case "up":
		m.searchResTable.MoveUp(1)
	case "down":
		m.searchResTable.MoveDown(1)
	}
	return m, nil
}

func (m model) handleReleaseResultInput(
	msg tea.KeyMsg,
) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.currentView = viewExporting
		m.buildExportData()
		if err := m.exportSongs(); err != nil {
			log.Printf("export error: %v", err)
		}
		m.currentView = viewSearch
	case "esc":
		m.currentView = viewSearch
	case "up":
		m.releaseDataTable.MoveUp(1)
	case "down":
		m.releaseDataTable.MoveDown(1)
	}
	return m, nil
}

func (m *model) searchRelease() tea.Cmd {
	return func() tea.Msg {
		resData, err := m.mb.SearchRelease(
			m.inputs[0].Value(),
			m.inputs[1].Value(),
			"12vinyl",
		)
		if err != nil {
			return errMsg{err}
		}
		return searchRes(resData)
	}
}

func (m *model) getReleaseData() tea.Cmd {
	return func() tea.Msg {
		id := m.searchRes.Releases[m.searchResTable.Cursor()].ID
		err := m.mb.GetReleaseData(id)
		if err != nil {
			return errMsg{err}
		}
		return m.mb.ReleaseData
	}
}
