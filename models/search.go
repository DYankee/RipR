package models

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Define styles
var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurStyle.Render("Submit"))

	windowStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1a1a1a")).
			Padding(1, 2).
			Margin(1)

	debugStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("240")).
			PaddingLeft(2).
			MarginTop(1)
)

type SearchModel struct {
	focusIndex int
	cursorMode cursor.Mode
	inputs     []textinput.Model
}

func InitialSearchModel() SearchModel {
	m := SearchModel{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 100
		switch i {
		case 0:
			t.SetWidth(20)
			t.Placeholder = "Artist"
			t.Focus()
		case 1:
			t.SetWidth(20)
			t.Placeholder = "Album"
		}
		m.inputs[i] = t
	}
	return m
}

func (m SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down", "enter":
			s := msg.String()

			// SUBMIT LOGIC: Instead of tea.Quit, send the message
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, func() tea.Msg {
					return SearchSubmittedMsg{
						Artist: m.inputs[0].Value(),
						Album:  m.inputs[1].Value(),
					}
				}
			}

			// Cycle focus
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
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *SearchModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m SearchModel) View() tea.View {
	var b strings.Builder
	var c *tea.Cursor

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View() + "\n")
		if m.inputs[i].Focused() {
			c = m.inputs[i].Cursor()
			if c != nil {
				c.Y += i
			}
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n%s\n", *button)
	v := tea.NewView(windowStyle.Render(b.String()))
	v.Cursor = c
	return v
}
