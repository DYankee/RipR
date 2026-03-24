package main

import (
	"log"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

func (m *model) buildReleaseTable(resData searchRes) {
	columns := []table.Column{
		{Title: "#", Width: 5},
		{Title: "Release Name", Width: 20},
		{Title: "Artist", Width: 10},
		{Title: "Country", Width: 10},
		{Title: "Year", Width: 5},
	}

	if len(resData.Releases) == 0 {
		log.Println("No search results found")
		return
	}

	rows := make([]table.Row, len(resData.Releases))
	for i, r := range resData.Releases {
		rows[i] = table.Row{
			strconv.Itoa(i),
			r.Title,
			r.ArtistCredit.NameCredits[0].Artist.Name,
			r.CountryCode,
			strconv.Itoa(r.Date.Year()),
		}
	}

	m.searchResTable = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
	)
}

func (m *model) buildReleaseResTable() {
	columns := []table.Column{
		{Title: "Track #", Width: 8},
		{Title: "Name", Width: 20},
		{Title: "Length", Width: 10},
	}

	var rows []table.Row
	for _, med := range m.releaseData.Mediums {
		for _, t := range med.Tracks {
			length := time.Millisecond * time.Duration(t.Length)
			rows = append(rows, table.Row{
				strconv.Itoa(t.Position),
				t.Recording.Title,
				length.String(),
			})
		}
	}

	m.releaseDataTable = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
	)
}
