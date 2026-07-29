package nav

import (
	mb "github.com/DYankee/RRipper/internal/musicbrainz"
)

// Types for navigating around the app

type GoBack struct{}
type ToSearch struct{}

type ToResults struct {
	Releases []mb.Release
}

type ToAlbum struct {
	Release *mb.Release
}

type ToQuantize struct {
	Release *mb.Release
}
