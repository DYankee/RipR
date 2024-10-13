package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	Internal "github.com/DYankee/RRipper/internal"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michiwend/gomusicbrainz"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type view int

const (
	SEARCH view = iota
	LOADING
	SEARCH_RESULT
	RELEASE_RESULT
)

func (v view) String() string {
	return [...]string{"Search", "Loading", "Search result", "Release result"}[v-1]
}

type searchRes gomusicbrainz.ReleaseSearchResponse
type releaseData gomusicbrainz.Release

type errMsg struct {
	err error
}

type songData struct {
	songName   string
	songLength float64
}

type sideData struct {
	clipInfo       Internal.ClipInfo
	songExportData []songData
	lengthMod      float64
}

func (sd *sideData) calcLengthMod() {

}

func (sd *sideData) getSideLength() float64 {
	// log.Println("getting side length")
	// log.Println(sd.sideData.End)
	// log.Println(sd.sideData.Start)
	len := sd.clipInfo.End - float64(sd.clipInfo.Start)
	// log.Println(len)
	return len
}

func (m *model) GetlengthMod() {
	var sideLength float64
	for k, v := range m.sideData {
		for _, v := range v.songExportData {
			sideLength += float64(v.songLength)
		}
		sideLength /= 1000
		log.Printf("Side length: %f", sideLength)
		log.Printf("Audacity Side length %f", m.sideData[k].getSideLength())

		dif := math.Abs(sideLength - m.sideData[k].getSideLength())
		log.Printf("Length difference: %f", dif)
		total := sideLength + m.sideData[k].getSideLength()
		log.Printf("Length total: %f", total)

		m.sideData[k].lengthMod = ((dif / total) / 2)
		log.Printf("Side Length mod: %f", m.sideData[k].lengthMod)

		sideLength = 0
	}
	//mod := (float64(k.Length) / 1000) - (float64(k.Length)/1000)*(m.sideData[x].lengthMod)

}

// Model and its functions
type model struct {
	mb               Internal.MusicBrainz
	audacity         Internal.Audacity
	searchRes        searchRes
	searchResTable   table.Model
	releaseData      releaseData
	releaseDataTable table.Model
	sideData         []sideData
	currentView      view
	inputs           []textinput.Model
	focusIndex       int
	focusIndexY      int
	Width            int
	Height           int
}

func New() *model {
	m := model{
		currentView: SEARCH,
		inputs:      make([]textinput.Model, 2),
	}
	m.mb.Init()
	m.audacity.Init()
	for !m.audacity.Status {
		println("connecting")
		m.audacity.Connect()
		time.Sleep(10000)
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
		{Title: "MB ID", Width: 10},
		{Title: "Release Name", Width: 30},
		{Title: "Artist", Width: 10},
		{Title: "County", Width: 10},
		{Title: "Year", Width: 30},
	}

	//build rows
	if len(resData.Releases) == 0 {
		panic(100)
	}
	rows := make([]table.Row, len(resData.Releases))
	for i, k := range resData.Releases {
		rows[i] = table.Row{strconv.Itoa(i), string(k.ID), k.Title, k.ArtistCredit.NameCredits[0].Artist.Name, k.CountryCode, strconv.Itoa(k.Date.Year())}
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
		resData := m.mb.ReleaseData
		return releaseData(resData)
	}
}

func (m *model) buildReleaseResTable(rd releaseData) {
	columns := []table.Column{
		{Title: "Track #", Width: 10},
		{Title: "Name", Width: 30},
		{Title: "length", Width: 30},
	}
	log.Println(rd)
	//build rows

	rows := make([]table.Row, 0)
	for _, k := range rd.Mediums {
		for _, k := range k.Tracks {
			length := time.Millisecond * time.Duration(k.Length)
			rows = append(rows, table.Row{k.Number, k.Recording.Title, length.String()})
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows))
	m.currentView = LOADING
	m.releaseDataTable = t
}

func (m *model) buildExportData() {
	data := m.audacity.GetClips()
	log.Println("Printing audacity side data")
	log.Println(data)

	var curSide, oldSide byte
	oldSide = 'A'
	m.sideData = append(m.sideData, sideData{})
	sideIdx := 0

	for _, medium := range m.releaseData.Mediums {
		for k2, t := range medium.Tracks {
			curSide = t.Number[0]
			log.Println(curSide)
			if curSide != oldSide {
				m.sideData = append(m.sideData, sideData{})
				sideIdx++
			}
			songName := fmt.Sprintf("%c%d - "+t.Recording.Title, t.Number[0], k2+1)
			m.sideData[sideIdx].songExportData = append(m.sideData[sideIdx].songExportData, songData{
				songLength: float64(t.Length),
				songName:   songName,
			})
			m.sideData[sideIdx].clipInfo = data[sideIdx]
			oldSide = curSide
		}
	}

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

	log.Println(m.sideData[1].songExportData[0].songLength)

	var offSet float64
	for _, sd := range m.sideData {
		log.Println(sd)
		log.Println(sd.clipInfo)

		offSet = float64(sd.clipInfo.Start)
		for _, s := range sd.songExportData {
			s.songLength = (s.songLength / 1000) + ((s.songLength / 1000) * sd.lengthMod)
			log.Printf("Exporting songs")
			log.Printf("offset: %f song length: %f", offSet, s.songLength)

			res := m.audacity.SelectRegion(offSet, offSet+s.songLength)
			log.Println("Select res:" + res)
			res = m.audacity.ExportAudio("./code/rripper/testdata/thriller", s.songName+".flac")
			log.Println("Export res:" + res)
			offSet += s.songLength
			log.Println(offSet)
		}
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
		m.buildReleaseTable(msg)
		m.searchRes = msg
		m.currentView = SEARCH_RESULT
		return m, cmd
	case releaseData:
		m.buildReleaseResTable(msg)
		m.releaseData = msg
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
			m.buildExportData()
			m.GetlengthMod()
			m.ExportSongs()
			m.currentView = SEARCH
		case "esc":
			m.currentView = SEARCH
		case "up", "down":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.releaseDataTable.MoveUp(1)
			} else {
				m.releaseDataTable.MoveDown(1)
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
		view = m.ResultView()
	case RELEASE_RESULT:
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

func (m *model) Header() string {
	headerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		PaddingBottom(m.Height / 3).
		Width(m.Width - 2)

	head := lipgloss.JoinVertical(lipgloss.Top,
		"Current view: "+m.currentView.String(),
		"Index X: "+strconv.Itoa(m.focusIndex),
		"Index Y: "+strconv.Itoa(m.focusIndex),
	)
	head = lipgloss.JoinHorizontal(lipgloss.Center,
		head,
		//"Artist: "+m.mb.ReleaseData.ArtistCredit.NameCredits[0].Artist.Name,
		//"Release: "+m.mb.ReleaseData.Title,
	)

	return headerStyle.Render(head)
}

func (m *model) Footer() string {
	footerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))
	footerStyle.PaddingTop(m.Height / 3).
		Width(m.Width - 2)
	return footerStyle.Render("test")
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
			m.Header(),
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.JoinVertical(
					lipgloss.Center,
					m.inputs[0].View(),
					m.inputs[1].View()),
				*button),
			m.Footer(),
		),
	)
}

func (m *model) ResultView() string {
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		m.searchResTable.View())
}

func (m *model) ReleaseView() string {
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		m.releaseDataTable.View())
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
