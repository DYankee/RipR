package main

import (
	"github.com/DYankee/godacity/clips"
	"github.com/michiwend/gomusicbrainz"
)

// view represents the current screen state.
type view int

const (
	viewSearch view = iota
	viewLoading
	viewSearchResult
	viewReleaseResult
	viewExporting
)

func (v view) String() string {
	return [...]string{
		"Search",
		"Loading",
		"Search result",
		"Release result",
		"Exporting",
	}[v]
}

// Message types for the Bubble Tea update loop.

type searchRes gomusicbrainz.ReleaseSearchResponse

type errMsg struct {
	err error
}

// cmds for changing current view
type ShowSearchMsg struct{}
type ShowResultsMsg struct{ Data searchRes }
type ShowExportMsg struct{ Release gomusicbrainz.Release }

// songData holds metadata for a single track to export.
type songData struct {
	songName     string
	songLength   float64 // in milliseconds
	songPosition int
}

// sideData represents one side/clip of the vinyl.
type sideData struct {
	clipInfo       clips.ClipInfo
	songExportData []songData
	sideLength     float64 // in seconds
	lengthMod      float64
}
