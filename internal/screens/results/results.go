// models/results.go
package results

import (
	"strconv"
	"strings"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	mb "github.com/DYankee/RRipper/internal/musicbrainz"
	"github.com/DYankee/RRipper/internal/nav"
	"github.com/DYankee/RRipper/internal/styles"
)

type Model struct {
	releases []mb.Release
	table    table.Model
}

func New(releases []mb.Release) Model {
	return Model{
		releases: releases,
		table:    buildTable(releases),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			idx := m.table.Cursor()
			if idx >= 0 && idx < len(m.releases) {
				r := m.releases[idx] // copy so the closure captures a stable value
				return m, func() tea.Msg {
					return nav.ToAlbum{Release: &r}
				}
			}
		case "esc":
			return m, func() tea.Msg { return nav.GoBack{} }
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Search Results"))
	b.WriteString("\n\n")

	if len(m.releases) == 0 {
		b.WriteString(styles.BlurStyle.Render("No releases found."))
	} else {
		// .Content extracts the rendered string from the sub-component's tea.View
		b.WriteString(m.table.View())
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render(
			"[↑/↓] Navigate  [enter] Select  [esc] Back",
		))
	}

	v := tea.NewView(styles.WindowStyle.Render(b.String()))
	return v
}

func buildTable(releases []mb.Release) table.Model {
	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Release", Width: 30},
		{Title: "Artist", Width: 20},
		{Title: "Date", Width: 12},
	}

	rows := make([]table.Row, len(releases))
	for i, r := range releases {
		rows[i] = table.Row{
			strconv.Itoa(i + 1),
			r.Title,
			r.Artist,
			r.Date,
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)
	t.SetStyles(tableStyles())
	return t
}

func tableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		Bold(true).
		Background(lipgloss.Color(styles.BackgroundColor))
	s.Selected = s.Selected.
		Background(lipgloss.Color(styles.BackgroundColor)).
		Bold(true)
	return s
}
