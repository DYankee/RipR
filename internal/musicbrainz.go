package Internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/michiwend/gomusicbrainz"
	"github.com/ryanuber/columnize"
)

type ReleaseQuerry struct {
	album  string
	artist string
	format string
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
	querry := fmt.Sprintf("release:%s AND artist:%s AND format:%s", q.album, q.artist, q.format)

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

func (m *MusicBrainz) GetQuerry() {
	reader := bufio.NewReader(os.Stdin)

	querryData := ReleaseQuerry{}
	// get album name
	fmt.Println("Please enter the album name")
	album, err := reader.ReadString('\n')
	if err != nil {
		fmt.Print(err)
	}
	querryData.album = strings.Replace(album, "\n", "", -1)

	fmt.Println("Please enter the artist name")
	artist, err := reader.ReadString('\n')
	if err != nil {
		fmt.Print(err)
	}
	querryData.artist = strings.Replace(artist, "\n", "", -1)

	fmt.Println("Please Select a release format")
	fmt.Println(`1: 12" Vinyl`)
	fmt.Println(`2: CD`)
	fmt.Println(`3: Digital Media`)

	res, err := reader.ReadByte()
	if err != nil {
		fmt.Println(err)
	}
	switch res {
	case '1':
		querryData.format = "12vinyl"
	case '2':
		querryData.format = "cd"
	case '3':
		querryData.format = "digital media"
	}

	fmt.Printf("Album: %s Artist: %s Format: %s \n", querryData.album, querryData.artist, querryData.format)

	m.ReleaseQuerrys = append(m.ReleaseQuerrys, querryData)
}

func (m *MusicBrainz) TestQuerry(album string, artist string, format string) {
	querryData := ReleaseQuerry{
		artist: artist,
		album:  album,
		format: format,
	}
	m.ReleaseQuerrys = append(m.ReleaseQuerrys, querryData)
}
