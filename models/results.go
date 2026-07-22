// models/results.go
package models

import (
	"log"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/bubbles/table"
	"github.com/michiwend/gomusicbrainz"
)

type ResultsModel struct {
	Releases      []*gomusicbrainz.Release
	ReleasesTable table.Model
}

func NewResultsModel(response gomusicbrainz.ReleaseSearchResponse) ResultsModel {
	return ResultsModel{
		Releases:      response.Releases,
		ReleasesTable: buildReleaseTable(response),
	}
}

func (m ResultsModel) Init() tea.Cmd {
	return nil
}

func (m ResultsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m ResultsModel) View() tea.View {
	var b strings.Builder
	b.WriteString("Search Results:\n\n")

	if len(m.Releases) == 0 {
		b.WriteString("No releases found.")
	} else {
		b.WriteString(m.ReleasesTable.View())
	}
	return tea.NewView(b.String())
}

// Helper functions
func buildReleaseTable(resData gomusicbrainz.ReleaseSearchResponse) table.Model {
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
	return t
}
