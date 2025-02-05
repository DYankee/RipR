package Internal

import (
	"fmt"

	"github.com/michiwend/gomusicbrainz"
	"github.com/ryanuber/columnize"
)

type ReleaseQuerry struct {
	Album  string
	Artist string
	Format string
}

type MusicBrainz struct {
	Client                *gomusicbrainz.WS2Client
	ReleaseSearchResponse gomusicbrainz.ReleaseSearchResponse
	ReleaseQuerrys        []ReleaseQuerry
	ReleaseData           gomusicbrainz.Release
}

func (m *MusicBrainz) DisplayReleaseData() {
	var id int
	output := []string{}
	output = append(output, "# | Posistion | Title | Length1 | Length2 | Artist")
	for _, release := range m.ReleaseData.Mediums[1].Tracks {
		id += 1
		formated := fmt.Sprintf("%d | %d | %s | %d | %d",
			id,
			release.Position,
			release.Recording.Title,
			release.Length,
			release.Recording.Length)
		output = append(output, formated)
	}
	fin := columnize.SimpleFormat(output)
	fmt.Println(fin)
}

func (m *MusicBrainz) Init() error {
	var err error
	m.Client, err = gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"RipR",
		"0.2.0-beta",
		"http://github.com/Dyankee/RRipper")
	if err != nil {
		return err
	}
	return nil
}

func (m *MusicBrainz) GetReleasesByFormat(resData *gomusicbrainz.ReleaseSearchResponse, format string) (res gomusicbrainz.ReleaseSearchResponse) {

	for i := 0; i < len(resData.Releases); i++ {
		if resData.Releases[i].Mediums[0].Format == format {
			res.Releases = append(res.Releases, resData.Releases[i])
		}
	}
	return res
}

func (m *MusicBrainz) DisplayReleaseRes(resData *gomusicbrainz.ReleaseSearchResponse) {
	var id int
	output := []string{}
	output = append(output, "# | Name | Artist | Format | Release Date | Country")
	for _, release := range resData.Releases {
		id += 1
		formated := fmt.Sprintf("%d | %s | %s | %s | %d | %s",
			id,
			release.Title,
			release.ArtistCredit.NameCredits[0].Artist.Name,
			release.Mediums[0].Format,
			release.Date.Year(),
			release.CountryCode)
		output = append(output, formated)
	}
	fin := columnize.SimpleFormat(output)
	fmt.Println(fin)
}

func (m *MusicBrainz) SearchRelease(artist string, release string, format string) (gomusicbrainz.ReleaseSearchResponse, error) {
	querry := fmt.Sprintf("release:%s AND artist:%s AND format:%s", release, artist, format)
	res, err := m.Client.SearchRelease(querry, -1, -1)
	return *res, err
}

func (m *MusicBrainz) GetReleaseData(id gomusicbrainz.MBID) error {
	data, err := m.Client.LookupRelease(id, "media+recordings")
	if err != nil {
		return err
	}
	m.ReleaseData = *data
	return nil
}

func (m *MusicBrainz) GetQuerry(album string, artist string) {
	querryData := ReleaseQuerry{
		Album:  album,
		Artist: artist,
		Format: "12vinyl",
	}
	m.ReleaseQuerrys = append(m.ReleaseQuerrys, querryData)
}

func (m *MusicBrainz) TestQuerry(album string, artist string, format string) {
	querryData := ReleaseQuerry{
		Artist: artist,
		Album:  album,
		Format: format,
	}
	m.ReleaseQuerrys = append(m.ReleaseQuerrys, querryData)
}
