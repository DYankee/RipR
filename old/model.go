package main

import (
	"fmt"
	"time"

	Internal "github.com/DYankee/RRipper/internal"
	"github.com/DYankee/godacity"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michiwend/gomusicbrainz"
)

type model struct {
	mb               Internal.MusicBrainz
	audacity         *godacity.Audacity
	searchRes        searchRes
	searchResTable   table.Model
	releaseData      gomusicbrainz.Release
	releaseDataTable table.Model
	sideData         []sideData
	currentView      view
	inputs           []textinput.Model
	focusIndex       int
	width            int
	height           int
}

func New() *model {
	m := model{
		currentView: viewSearch,
		inputs:      make([]textinput.Model, 3),
	}

	if err := m.mb.Init(); err != nil {
		panic(fmt.Sprintf("failed to init MusicBrainz: %v", err))
	}

	var err error
	m.audacity, err = godacity.NewAudacity(&godacity.Config{
		AutoStart:    true,
		StartTimeout: 15 * time.Second,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to Audacity: %v", err))
	}

	clips, err := m.audacity.Clips.GetClips()
	if err != nil {
		panic(fmt.Sprintf("failed to get clips: %v", err))
	}
	for _, c := range clips {
		m.sideData = append(m.sideData, sideData{clipInfo: c})
	}

	m.initInputs()
	return &m
}

func (m *model) initInputs() {
	placeholders := []string{"Artist", "Release", "Output Directory"}
	for i := range m.inputs {
		t := textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		t.Placeholder = placeholders[i]

		if i == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}

		m.inputs[i] = t
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
