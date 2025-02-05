package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	Audacity "github.com/DYankee/RRipper/internal/audacity"
	MusicBrainz "github.com/DYankee/RRipper/internal/musicbrainz"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michiwend/gomusicbrainz"
)

var (
	DefaultBorderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

// List of view states
type view int

const (
	SEARCH view = iota
	LOADING
	SEARCH_RESULT
	RELEASE_RESULT
	EXPORTING
)

func (v view) String() string {
	return [...]string{"Search", "Loading", "Search result", "Release result", "Exporting"}[v]
}

type searchRes gomusicbrainz.ReleaseSearchResponse

type errMsg struct {
	err error
}

type songData struct {
	songName     string
	songLength   float64
	songPosition int
}

type sideData struct {
	clipInfo       Audacity.ClipInfo
	songExportData []songData
	sideEnd        int
	sideLength     float64
	lengthMod      float64
}

func (sd *sideData) calcLengthMod() {
	if sd.songExportData != nil {
		for _, v := range sd.songExportData {
			sd.sideLength += float64(v.songLength)
		}
		log.Printf("Side length: %f", sd.sideLength)
		sd.sideLength /= 1000
		log.Printf("converted Side length: %f", sd.sideLength)
		log.Printf("Audacity Side length %f", sd.clipInfo.GetClipLength())

		dif := math.Abs(sd.sideLength - sd.clipInfo.GetClipLength())
		log.Printf("Length difference: %f", dif)
		total := sd.sideLength + sd.clipInfo.GetClipLength()
		log.Printf("Length total: %f", total)

		sd.lengthMod = ((dif / total) / 2)
		log.Printf("Length mod: %f", total)

	} else {
		log.Fatal("No song export data for side")
	}
}

// Model and its functions
type model struct {
	mb               MusicBrainz.MusicBrainz
	audacity         Audacity.Audacity
	searchRes        searchRes
	searchResTable   table.Model
	releaseData      gomusicbrainz.Release
	releaseDataTable table.Model
	sideData         []sideData
	sideIdx          int
	currentView      view
	inputs           []textinput.Model
	focusIndex       int
	Width            int
	Height           int
}

func New() *model {
	m := model{
		currentView: SEARCH,
		inputs:      make([]textinput.Model, 3),
		sideIdx:     0,
	}
	m.mb.Init()
	m.audacity.Init()
	for !m.audacity.Status {
		println("connecting")
		m.audacity.Connect()
		time.Sleep(10000)
	}
	data := m.audacity.GetClips()
	for _, v := range data {
		m.sideData = append(m.sideData, sideData{
			clipInfo: v,
		})
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
		case 2:
			t.Placeholder = "Output Directory"
		}

		m.inputs[i] = t
	}

	return &m
}

func (m *model) searchRelease() tea.Cmd {
	return func() tea.Msg {
		resData, err := m.mb.SearchRelease(m.inputs[0].Value(), m.inputs[1].Value(), "12vinyl")
		if err != nil {
			return errMsg{err}
		}
		return searchRes(resData)
	}
}

func (m *model) buildReleaseTable(resData searchRes) {
	columns := []table.Column{
		{Title: "#", Width: 5},
		{Title: "Release Name", Width: 20},
		{Title: "Artist", Width: 10},
		{Title: "County", Width: 10},
		{Title: "Year", Width: 5},
	}

	//build rows
	if len(resData.Releases) == 0 {
		log.Fatal("No results")
	}
	rows := make([]table.Row, len(resData.Releases))
	for i, k := range resData.Releases {
		rows[i] = table.Row{strconv.Itoa(i), k.Title, k.ArtistCredit.NameCredits[0].Artist.Name, k.CountryCode, strconv.Itoa(k.Date.Year())}
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows))
	m.searchResTable = t
}

func (m *model) GetReleaseData() tea.Cmd {
	return func() tea.Msg {
		log.Println(m.searchResTable.Cursor())
		log.Println(m.searchRes)
		//log.Println(m.mb.ReleaseSearchResponse.Releases)
		id := m.searchRes.Releases[m.searchResTable.Cursor()].ID
		err := m.mb.GetReleaseData(id)
		if err != nil {
			return errMsg{err}
		}
		return m.mb.ReleaseData
	}
}

func (m *model) buildReleaseResTable() {
	columns := []table.Column{
		{Title: "Track #", Width: 8},
		{Title: "Name", Width: 20},
		{Title: "length", Width: 10},
	}
	//build rows

	rows := make([]table.Row, 0)
	for _, med := range m.releaseData.Mediums {
		for _, t := range med.Tracks {
			length := time.Millisecond * time.Duration(t.Length)
			rows = append(rows, table.Row{strconv.Itoa(t.Position), t.Recording.Title, length.String()})
		}
	}
	m.releaseDataTable = table.New(
		table.WithColumns(columns),
		table.WithRows(rows))
}

func (m *model) buildExportData() {
	for _, Mediums := range m.releaseData.Mediums {
		for _, t := range Mediums.Tracks {
			m.sideData[0].songExportData = append(m.sideData[0].songExportData, songData{
				songName:     t.Recording.Title,
				songLength:   float64(t.Recording.Length),
				songPosition: t.Position,
			})
		}
	}
	m.sideData[0].calcLengthMod()
}

func (m *model) ExportSongs() {
	for _, sd := range m.sideData {
		for _, s := range sd.songExportData {
			log.Printf("Song length %f", s.songLength)
			s.songLength = s.songLength + (s.songLength * sd.lengthMod)
			log.Printf("Song length after mod %f", s.songLength)
			s.songLength = s.songLength / 1000
			log.Printf("Song length converted to correct time %f", s.songLength)
		}
	}
	var offSet float64
	for _, sd := range m.sideData {
		log.Println(sd)
		log.Println(sd.clipInfo)
		offSet = float64(sd.clipInfo.Start)
		for _, s := range sd.songExportData {
			s.songLength = ((s.songLength / 1000) + ((s.songLength * sd.lengthMod) / 1000) + 1)
			log.Printf("Exporting songs")
			log.Printf("offset: %f song length: %f", offSet, s.songLength)
			res := m.audacity.SelectRegion(offSet, offSet+s.songLength)
			log.Println("Select res:" + res)
			os.Mkdir(m.inputs[2].Value(), 0700)
			wd, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}
			res = m.audacity.ExportAudio(wd+"/"+m.inputs[2].Value(), strconv.Itoa(s.songPosition)+"_"+s.songName+".flac")
			log.Println("Export res:" + res)
			offSet += s.songLength
			log.Println(offSet)
		}
	}
}

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
		m.buildReleaseTable(msg)
		m.searchRes = msg
		m.focusIndex = 0
		m.currentView = SEARCH_RESULT
		return m, cmd
	case gomusicbrainz.Release:
		m.releaseData = msg
		m.buildReleaseResTable()
		m.currentView = RELEASE_RESULT
		return m, cmd
	}
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
	case SEARCH:
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
				m.currentView = LOADING
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
	case SEARCH_RESULT:
		switch msg.String() {
		case "enter":
			return m, m.GetReleaseData()
		case "esc":
			m.currentView = SEARCH
		case "up", "down":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.searchResTable.MoveUp(1)
			} else {
				m.searchResTable.MoveDown(1)
			}

		}
	case RELEASE_RESULT:
		switch msg.String() {
		case "enter":
			m.currentView = EXPORTING
			m.buildExportData()
			m.ExportSongs()
			m.currentView = SEARCH
		case "esc":
			m.currentView = SEARCH
		case "up", "down":
			s := msg.String()

			// Cycle indexes
			if s == "down" {
				m.releaseDataTable.MoveDown(1)
			} else {
				m.releaseDataTable.MoveUp(1)
			}

		}
	}
	if m.currentView == SEARCH {
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
	// log.Println("batch cmd abt to be run")
	return tea.Batch(cmds...)
}

// ------------------------------------------------------------------------------------

// view and helper functions
func (m model) View() string {
	var view string
	switch m.currentView {
	case SEARCH:
		view = m.SearchView()
	case LOADING:
		view = m.LoadingView()
	case SEARCH_RESULT:
		view = m.SearchResultView()
	case RELEASE_RESULT:
		view = m.ReleaseView()
	case EXPORTING:
		view = m.ExportingView()
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

func (m *model) Header() string {
	headerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(m.Width - 2).Align(lipgloss.Center)
	head := lipgloss.JoinHorizontal(lipgloss.Center,
		"| Current view: "+m.currentView.String()+" |",
		"| Index X: "+strconv.Itoa(m.focusIndex)+" |",
	)
	switch m.currentView {
	case SEARCH_RESULT:
		head = lipgloss.JoinVertical(lipgloss.Center,
			head,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				"| Current Search artist: "+m.inputs[0].Value()+" |",
				"| Current Search release: "+m.inputs[1].Value()+" |",
			),
		)
	case RELEASE_RESULT:
		head = lipgloss.JoinVertical(lipgloss.Center,
			head,
			"| Current Release: "+m.releaseData.Title+" | "+m.releaseData.Disambiguation+" |",
			"| Cursor pos: "+strconv.Itoa(m.releaseDataTable.Cursor()),
		)
	}
	return headerStyle.Render(head)
}

func (m *model) SearchView() string {
	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	headerH := lipgloss.Height(m.Header())
	log.Println(headerH)
	bodystyle := lipgloss.NewStyle().
		Width(m.Width - 2).
		PaddingTop((m.Height / 2) - (headerH)).
		PaddingBottom((m.Height / 2) - (headerH + 2)).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		AlignHorizontal(lipgloss.Center)

	body := lipgloss.JoinVertical(
		lipgloss.Center,
		m.inputs[0].View(),
		m.inputs[1].View(),
		m.inputs[2].View(),
		*button)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.Header(),
			bodystyle.Render(body),
		),
	)
}

func (m *model) SearchResultView() string {
	headerH := lipgloss.Height(m.Header())
	bodystyle := lipgloss.NewStyle().
		Width(m.Width - 2).
		PaddingTop((m.Height / 2) - (headerH + 20)).
		PaddingBottom((m.Height / 2) - (headerH)).
		MarginTop(0).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		AlignHorizontal(lipgloss.Center)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.Header(),
			bodystyle.Render(m.searchResTable.View()),
		))
}

func (m *model) ReleaseView() string {

	button := &blurredButton
	if m.focusIndex == 1 {
		button = &focusedButton
	}
	headerH := lipgloss.Height(m.Header())
	bodystyle := lipgloss.NewStyle().
		Width(m.Width - 2).
		PaddingTop((m.Height / 2) - (headerH + 20)).
		PaddingBottom((m.Height / 2) - (headerH)).
		MarginTop(0).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		AlignHorizontal(lipgloss.Center)

	body := lipgloss.JoinVertical(
		lipgloss.Center,
		m.releaseDataTable.View(),
		*button,
	)
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.Header(),
			bodystyle.Render(body),
		),
	)
}

func (m *model) ExportingView() string {
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, "Exporting")
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
