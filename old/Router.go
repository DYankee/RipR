package main

import (
	Internal "github.com/DYankee/RRipper/internal"
	"github.com/DYankee/godacity"
	tea "github.com/charmbracelet/bubbletea"
)

type Router struct {
	activeModel tea.Model
	// Shared resources
	mb       Internal.MusicBrainz
	audacity *godacity.Audacity
}

func (r Router) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 1. Handle Global Interrupts (Ctrl+C)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "ctrl+c" {
			return r, tea.Quit
		}
	}

	// 2. Handle Navigation Messages (The "Handlers")
	// Instead of a switch on an enum, we switch on the MESSAGE TYPE
	switch msg := msg.(type) {
	case ShowSearchMsg:
		r.activeModel = NewSearchModel()
		return r, r.activeModel.Init()

	case ShowResultsMsg:
		r.activeModel = NewResultsModel(msg.Data)
		return r, r.activeModel.Init()
	}

	// 3. Delegate EVERYTHING else to the active sub-model
	var cmd tea.Cmd
	r.activeModel, cmd = r.activeModel.Update(msg)
	return r, cmd
}

func (r Router) View() string {
	// Just render whatever the active sub-model is
	return r.activeModel.View()
}
