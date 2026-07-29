package musicbrainz

import (
	"fmt"

	"github.com/michiwend/gomusicbrainz"
)

type MusicBrainz struct {
	Client *gomusicbrainz.WS2Client
}

func NewClient() (*MusicBrainz, error) {
	Client, err := gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"RipR",
		"0.2.0-beta",
		"http://github.com/Dyankee/RRipper",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to musicbrainz client: %w", err)
	}

	return &MusicBrainz{Client: Client}, nil
}

func (m *MusicBrainz) SearchRelease(artist, release, format string) ([]Release, error) {
	query := fmt.Sprintf(
		"release:%s AND artist:%s AND format:%s",
		release, artist, format,
	)
	res, err := m.Client.SearchRelease(query, -1, -1)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	return mapReleases(res.Releases), nil
}

// map release data from musicbrainz to internal data structure
func mapReleases(releases []*gomusicbrainz.Release) []Release {
	out := make([]Release, 0, len(releases))
	for _, r := range releases {
		out = append(out, Release{
			ID:     string(r.ID),
			Title:  r.Title,
			Artist: artistCredit(r),
			Date:   r.Date.String(),
			Tracks: mapTracks(r),
		})
	}
	return out
}

// Extract the artist credit from release data
func artistCredit(r *gomusicbrainz.Release) string {
	if len(r.ArtistCredit.NameCredits) == 0 {
		return ""
	}
	return r.ArtistCredit.NameCredits[0].Artist.Name
}

// map track data from musicbrainz to internal data structure
func mapTracks(r *gomusicbrainz.Release) []Track {
	var tracks []Track
	for _, medium := range r.Mediums {
		for _, t := range medium.Tracks {
			tracks = append(tracks, Track{
				Position: t.Position,
				Title:    t.Recording.Title,
				Length:   t.Recording.Length,
			})
		}
	}
	return tracks
}
