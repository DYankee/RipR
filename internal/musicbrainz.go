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
	Client                 *gomusicbrainz.WS2Client
	ReleaseSearchResponses []*gomusicbrainz.ReleaseSearchResponse
	ReleaseQuerrys         []ReleaseQuerry
	ReleaseData            *gomusicbrainz.Release
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

func (m *MusicBrainz) ChooseRelease() {
	fmt.Println("Enter number of chosen release")
	var res int
	_, err := fmt.Scan(&res)
	if err != nil {
		fmt.Println(err)
	}

	m.GetReleaseData(res)
}

func (m *MusicBrainz) Init() error {
	var err error
	m.Client, err = gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"RRipper",
		"0.0.2-beta",
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

func (m *MusicBrainz) SearchRelease(q *ReleaseQuerry) error {
	querry := fmt.Sprintf("release:%s AND artist:%s AND format:%s", q.Album, q.Artist, q.Format)

	res, err := m.Client.SearchRelease(querry, -1, -1)
	if err != nil {
		return err
	}
	m.ReleaseSearchResponses = append(m.ReleaseSearchResponses, res)
	return nil
}

func (m *MusicBrainz) GetReleaseData(i int) {
	data, err := m.Client.LookupRelease(m.ReleaseSearchResponses[0].Releases[i].ID, "media+recordings")
	if err != nil {
		fmt.Println(err)
	}
	m.ReleaseData = data
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
