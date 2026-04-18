// models/results.go
package models

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/michiwend/gomusicbrainz"
)

type ResultsModel struct {
	Releases []gomusicbrainz.Release
}

func NewResultsModel(res gomusicbrainz.ReleaseSearchResponse) ResultsModel {
	return ResultsModel{
		Releases: res.Releases,
	}
}

func (m ResultsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m ResultsModel) View() tea.View {
	var b strings.Builder
	b.WriteString("Search Results:\n\n")

	if len(m.Releases) == 0 {
		b.WriteString("No releases found.")
	}

	for i, r := range m.Releases {
		if i > 10 {
			break
		} // Limit display
		b.WriteString(fmt.Sprintf("- %s (%s) [%s]\n", r.Title, r.Date.String(), r.ID))
	}

	return tea.NewView(b.String())
}
