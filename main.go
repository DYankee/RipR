package main

import (
	"log"
	"strconv"

	Internal "github.com/DYankee/RRipper/internal"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
	table       lipgloss.Style
}

var tableStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

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
	querry       Internal.ReleaseQuerry
	searchRes    table.Model
	answerField  textinput.Model
	audacity     Internal.Audacity
	musicBrainz  Internal.MusicBrainz
}

func New(searchfields []Question) *model {
	Styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Focus()
	mb := Internal.MusicBrainz{}
	mb.Init()
	return &model{
		currentView:  "main",
		searchfields: searchfields,
		answerField:  answerField,
		Styles:       Styles,
		querry:       Internal.ReleaseQuerry{},
		musicBrainz:  mb,
		audacity:     Internal.Audacity{},
	}
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
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
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
				m.answerField.SetValue("")
				log.Printf("Question: %s, Answer: %s", current.question, current.answer)
				m.Next()
				return m, nil
			}
		case "result":
			switch msg.String() {
			case "enter":
				m.currentView = "main"
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
		m.querry.Artist = m.searchfields[0].answer
		m.querry.Album = m.searchfields[1].answer
		m.querry.Format = "12vinyl"
		m.currentView = "result"
		log.Printf("Artist: %s, Release: %s, Format: %s", m.querry.Artist, m.querry.Album, m.querry.Format)
		m.musicBrainz.SearchRelease(&m.querry)
		m.musicBrainz.GetReleaseData(0)
		m.FormatTable()
	}
}

func (m *model) FormatTable() {
	resData := m.musicBrainz.ReleaseData.Mediums[0].Tracks
	colums := []table.Column{
		{Title: "#", Width: 5},
		{Title: "Track", Width: 30},
		{Title: "Length", Width: 10},
	}

	//build rows
	rows := []table.Row{
		{strconv.Itoa(resData[0].Position), resData[0].Recording.Title, strconv.Itoa(resData[0].Length)},
		{strconv.Itoa(resData[1].Position), resData[1].Recording.Title, strconv.Itoa(resData[1].Length)},
		{strconv.Itoa(resData[2].Position), resData[2].Recording.Title, strconv.Itoa(resData[2].Length)},
		{strconv.Itoa(resData[3].Position), resData[3].Recording.Title, strconv.Itoa(resData[3].Length)},
	}
	t := table.New(
		table.WithColumns(colums),
		table.WithRows(rows),
	)
	m.searchRes = t
	log.Printf("%s, %s", rows[0][0], rows[0][1])
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
	case "result":
		return tableStyle.Render(m.searchRes.View()) + "\n"
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
