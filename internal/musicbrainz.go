package Internal

import (
	"fmt"

	"github.com/michiwend/gomusicbrainz"
)

type MusicBrainz struct {
	Client *gomusicbrainz.WS2Client
}

func NewClient() *MusicBrainz {
	Client, err := gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"RipR",
		"0.2.0-beta",
		"http://github.com/Dyankee/RRipper",
	)
	if err != nil {
		println(err)
	}

	return &MusicBrainz{Client: Client}
}

func (m *MusicBrainz) SearchRelease(
	artist, release, format string,
) (gomusicbrainz.ReleaseSearchResponse, error) {
	query := fmt.Sprintf(
		"release:%s AND artist:%s AND format:%s",
		release, artist, format,
	)
	res, err := m.Client.SearchRelease(query, -1, -1)
	if err != nil {
		return gomusicbrainz.ReleaseSearchResponse{}, err
	}
	return *res, nil
}
