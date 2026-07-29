package main

import (
	// Std lib
	"fmt"
	"os"

	// External
	tea "charm.land/bubbletea/v2"

	//Internal
	"github.com/DYankee/RRipper/internal/app"
	mb "github.com/DYankee/RRipper/internal/musicbrainz"
)

func main() {
	client, err := mb.NewClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(app.New(client))
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
