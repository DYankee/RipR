package Internal

import (
	"fmt"

	"github.com/michiwend/gomusicbrainz"
	"github.com/ryanuber/columnize"
)

type ReleaseQuery struct {
	Album  string
	Artist string
	Format string
}

type MusicBrainz struct {
	Client         *gomusicbrainz.WS2Client
	ReleaseData    gomusicbrainz.Release
	ReleaseQueries []ReleaseQuery
}

func (m *MusicBrainz) Init() error {
	var err error
	m.Client, err = gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"RipR",
		"0.2.0-beta",
		"http://github.com/Dyankee/RRipper",
	)
	return err
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

func (m *MusicBrainz) GetReleaseData(id gomusicbrainz.MBID) error {
	data, err := m.Client.LookupRelease(id, "media+recordings")
	if err != nil {
		return err
	}
	m.ReleaseData = *data
	return nil
}

func (m *MusicBrainz) GetReleasesByFormat(
	resData *gomusicbrainz.ReleaseSearchResponse,
	format string,
) gomusicbrainz.ReleaseSearchResponse {
	var res gomusicbrainz.ReleaseSearchResponse
	for _, r := range resData.Releases {
		if len(r.Mediums) > 0 && r.Mediums[0].Format == format {
			res.Releases = append(res.Releases, r)
		}
	}
	return res
}

func (m *MusicBrainz) DisplayReleaseRes(
	resData *gomusicbrainz.ReleaseSearchResponse,
) {
	output := []string{
		"# | Name | Artist | Format | Release Date | Country",
	}
	for i, r := range resData.Releases {
		output = append(output, fmt.Sprintf(
			"%d | %s | %s | %s | %d | %s",
			i+1, r.Title,
			r.ArtistCredit.NameCredits[0].Artist.Name,
			r.Mediums[0].Format,
			r.Date.Year(), r.CountryCode,
		))
	}
	fmt.Println(columnize.SimpleFormat(output))
}

func (m *MusicBrainz) DisplayReleaseData() {
	output := []string{
		"# | Position | Title | Length | Recording Length",
	}
	for i, t := range m.ReleaseData.Mediums[0].Tracks {
		output = append(output, fmt.Sprintf(
			"%d | %d | %s | %d | %d",
			i+1, t.Position, t.Recording.Title,
			t.Length, t.Recording.Length,
		))
	}
	fmt.Println(columnize.SimpleFormat(output))
}
