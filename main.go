package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type Question struct {
	question string
	answer   string
}

func NewQuestion(q string) Question {
	return Question{question: q}
}

type model struct {
	currentView  string
	Styles       *Styles
	width        int
	height       int
	searchIndex  int
	searchfields []Question
	answerField  textinput.Model
}

func New(searchfields []Question) *model {
	Styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Focus()
	return &model{
		currentView:  "main",
		searchfields: searchfields,
		answerField:  answerField,
		Styles:       Styles}
}

func (m model) Init() tea.Cmd {
	return nil
}

// update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.searchfields[m.searchIndex]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch m.currentView {
		case "main":
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				m.currentView = "input"
				return m, nil
			}
		case "input":
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				current.answer = m.answerField.Value()
				log.Printf("Question: %s, Answer: %s", current.question, current.answer)
				m.Next()
				m.answerField.SetValue("")
				return m, nil
			}
		}
	}
	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

// update helper functions
func (m *model) Next() {
	if m.searchIndex < len(m.searchfields)-1 {
		m.searchIndex++
	} else {
		m.searchIndex = 0
		m.currentView = "main"
	}
}

// view function
func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}
	switch m.currentView {
	case "main":
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,

			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.JoinHorizontal(
					lipgloss.Center,
					m.searchfields[0].question,
					m.searchfields[0].answer,
				),
				lipgloss.JoinHorizontal(
					lipgloss.Center,
					m.searchfields[1].question,
					m.searchfields[1].answer,
				),
			),
		)
	case "input":
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,

			lipgloss.JoinVertical(
				lipgloss.Center,
				m.searchfields[m.searchIndex].question,
				m.Styles.InputField.Render(m.answerField.View()),
			),
		)
	}
	return "error"
}

// view helper func

func main() {
	questions := []Question{
		NewQuestion("Artist Name: "),
		NewQuestion("Release Name: "),
	}
	m := New(questions)

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
