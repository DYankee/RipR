package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

// Model and its functions
type model struct {
	currentView string
	inputs      []textinput.Model
	focusIndex  int
	Width       int
	Height      int
}

func New() *model {
	m := model{
		currentView: "Search",
		inputs:      make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Artist"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Release"
		}

		m.inputs[i] = t
	}

	return &m
}

func (m model) Init() tea.Cmd {
	return nil
}

//-------------------------------------------------------------------

// Update and helper functions
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.inputHandler(msg)
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	// Handle character input and blinking
	var cmd tea.Cmd
	return m, cmd
}

func (m *model) inputHandler(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.currentView {
	case "Search":
		switch msg.String() {

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.nextView()
				log.Printf("Artist: %s Release: %s",
					m.inputs[0].Value(),
					m.inputs[1].Value(),
				)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}

	// result screan controls
	case "Result":
		switch msg.String() {
		case "enter":
			m.nextView()
		}
	}
	if m.currentView == "Search" {
		cmd = m.updateInputs(msg)
	}

	return m, cmd
}

func (m *model) nextView() {
	switch m.currentView {
	case "Search":
		m.currentView = "Result"
	case "Result":
		m.currentView = "Search"
	}
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// ------------------------------------------------------------------------------------

// view and helper functions
func (m model) View() string {
	var view string
	switch m.currentView {
	case "Search":
		view = m.formatSearchView()
	case "Result":
		view = m.formatResultView()
	}
	return view
}
func (m *model) formatSearchView() string {

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Center,
				m.inputs[0].View(),
				m.inputs[1].View()),
			*button),
	)
}

func (m *model) formatResultView() string {
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		fmt.Sprintf("Artist: %s Release: %s", m.inputs[0].Value(), m.inputs[1].Value()),
	)
}

func main() {
	m := New()

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
