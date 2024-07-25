package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	Internal "github.com/DYankee/RRipper/internal"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type errMsg struct {
	err error
}

type songData struct {
	songName   string
	songLength float64
}

// Model and its functions
type model struct {
	mb          Internal.MusicBrainz
	audacity    Internal.Audacity
	searchRes   table.Model
	releaseData table.Model
	ExportData  []songData
	TrackInfo   Internal.TrackInfo
	currentView string
	inputs      []textinput.Model
	focusIndex  int
	focusIndexY int
	Width       int
	Height      int
}

func New() *model {
	m := model{
		currentView: "Search",
		inputs:      make([]textinput.Model, 2),
	}
	m.mb.Init()
	m.audacity.Connect()
	m.TrackInfo = m.audacity.GetInfo()
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

type searchRes table.Model
type releaseData table.Model

func (m *model) searchRelease() tea.Cmd {
	return func() tea.Msg {
		err := m.mb.SearchRelease(m.inputs[0].Value(), m.inputs[1].Value(), "12vinyl")
		if err != nil {
			return errMsg{err}
		}
		resData := m.mb.ReleaseSearchResponses

		colums := []table.Column{
			{Title: "#", Width: 5},
			{Title: "MB ID", Width: 10},
			{Title: "Release Name", Width: 30},
			{Title: "Artist", Width: 10},
			{Title: "County", Width: 10},
			{Title: "Year", Width: 30},
		}

		//build rows
		rows := make([]table.Row, len(resData.Releases))
		for i, k := range resData.Releases {
			rows[i] = table.Row{strconv.Itoa(i), string(k.ID), k.Title, k.ArtistCredit.NameCredits[0].Artist.Name, k.CountryCode, strconv.Itoa(k.Date.Year())}
		}
		t := table.New(
			table.WithColumns(colums),
			table.WithRows(rows))
		m.searchRes = t
		return searchRes(t)
	}
}

func (m *model) GetReleaseData() tea.Cmd {
	return func() tea.Msg {
		log.Println(m.searchRes.Cursor())
		id := m.mb.ReleaseSearchResponses.Releases[m.searchRes.Cursor()].ID
		err := m.mb.GetReleaseData(id)
		if err != nil {
			return errMsg{err}
		}
		resData := m.mb.ReleaseData
		colums := []table.Column{
			{Title: "Track #", Width: 10},
			{Title: "Name", Width: 30},
			{Title: "length", Width: 30},
		}
		log.Println(resData)
		//build rows

		rows := make([]table.Row, 0)
		for _, k := range resData.Mediums {
			for _, k := range k.Tracks {
				length := time.Millisecond * time.Duration(k.Length)
				rows = append(rows, table.Row{k.Number, k.Recording.Title, length.String()})
			}
		}

		t := table.New(
			table.WithColumns(colums),
			table.WithRows(rows))
		m.currentView = "Loading"
		return releaseData(t)
	}
}

func (m *model) GetlengthMod() (lengthMod float64) {
	m.TrackInfo = m.audacity.GetInfo()
	var releaseLength float64
	for _, v := range m.mb.ReleaseData.Mediums {
		for _, v := range v.Tracks {
			releaseLength += float64(v.Length) / 1000
		}
	}
	log.Println(m.TrackInfo)
	log.Printf("Release length: %f", releaseLength)
	lengthMod = (releaseLength - m.TrackInfo.End) / ((releaseLength + m.TrackInfo.End) / 2)
	log.Printf("Length mod: %f", lengthMod)
	lengthMod *= .99
	log.Printf("Length mod: %f", lengthMod)
	return lengthMod
}

func (m *model) buildExportData() {
	lengthMod := m.GetlengthMod()
	data := make([]songData, 0)
	for _, k := range m.mb.ReleaseData.Mediums {
		for i, k := range k.Tracks {
			songName := strings.Replace(k.Recording.Title, " ", "-", -1)
			songName = fmt.Sprintf("%d-"+songName, i+1)
			data = append(data, songData{
				songLength: (float64(k.Length) / 1000) - (float64(k.Length)/1000)*lengthMod,
				songName:   songName,
			})
		}
	}
	log.Println(data)
	m.ExportData = data
}

func (m *model) ExportSongs() {

	var offSet float64
	offSet = 0.00
	for _, v := range m.ExportData {
		res := m.audacity.SelectRegion(offSet, offSet+v.songLength)
		log.Println("Select res:" + res)
		res = m.audacity.ExportAudio("./code/rripper/testdata", v.songName+".flac")
		log.Println("Export res:" + res)
		offSet += v.songLength
		log.Println(offSet)
	}
}

//func (m *model) ExportSongs() {
//
// }

func (m model) Init() tea.Cmd {
	return nil
}

//-------------------------------------------------------------------

// Update and helper functions
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case errMsg:
		m.errorHandler(msg)
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, cmd
	case tea.KeyMsg:
		return m.inputHandler(msg)
	case searchRes:
		m.currentView = "SearchResult"
		m.searchRes = table.Model(msg)
		return m, cmd
	case releaseData:
		m.currentView = "ReleaseResult"
		m.releaseData = table.Model(msg)
		return m, cmd
	}

	// Handle character input and blinking
	return m, cmd
}

func (m *model) errorHandler(msg errMsg) {
	log.Println(msg.err)
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
				log.Printf("Artist: %s Release: %s",
					m.inputs[0].Value(),
					m.inputs[1].Value(),
				)
				m.currentView = "Loading"
				return m, m.searchRelease()
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
		}

	// result screen controls
	case "SearchResult":
		switch msg.String() {
		case "enter":
			return m, m.GetReleaseData()
		case "esc":
			m.currentView = "Search"
		case "up", "down":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.searchRes.MoveUp(1)
			} else {
				m.searchRes.MoveDown(1)
			}

		}
	case "ReleaseResult":
		switch msg.String() {
		case "enter":
			m.buildExportData()
			m.ExportSongs()
			m.currentView = "Search"
		case "esc":
			m.currentView = "Search"
		case "up", "down":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.releaseData.MoveUp(1)
			} else {
				m.releaseData.MoveDown(1)
			}
		}
	}
	if m.currentView == "Search" {
		cmd = m.updateInputs(msg)
	}

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
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
	log.Println("batch cmd abt to be run")
	return tea.Batch(cmds...)
}

// ------------------------------------------------------------------------------------

// view and helper functions
func (m model) View() string {
	var view string
	switch m.currentView {
	case "Search":
		view = m.SearchView()
	case "Loading":
		view = m.LoadingView()
	case "SearchResult":
		view = m.ResultView()
	case "ReleaseResult":
		view = m.ReleaseView()
	}
	return view
}

func (m *model) LoadingView() string {
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		"loading",
	)
}

func (m *model) SearchView() string {
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

func (m *model) ResultView() string {
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		m.searchRes.View())
}

func (m *model) ReleaseView() string {
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		m.releaseData.View())
}

//-------------------------------------------------------------------------------------------------------------------------------------

// Main function
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
