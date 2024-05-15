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
}

type MusicBrainz struct {
	client                 *gomusicbrainz.WS2Client
	ReleaseSearchResponses []*gomusicbrainz.ReleaseSearchResponse
	ReleaseQuerrys         []ReleaseQuerry
}

func (m *MusicBrainz) Init() error {
	var err error
	m.client, err = gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"RRipper",
		"0.0.1-beta",
		"http://github.com/Dyankee/RRipper")
	if err != nil {
		return err
	}
	return nil
}

func (m *MusicBrainz) DisplayReleaseRes(resData *gomusicbrainz.ReleaseSearchResponse) {
	var id int
	output := []string{}
	output = append(output, "# | Name | Artist | Format | Release Date | Country")
	for _, release := range resData.Releases {
		id += 1
		if release.Mediums[0].Format == `12" Vinyl` {
			formated := fmt.Sprintf("%d | %s | %s | %s | %d | %s",
				id,
				release.Title,
				release.ArtistCredit.NameCredits[0].Artist.Name,
				release.Mediums[0].Format,
				release.Date.Year(),
				release.CountryCode)
			output = append(output, formated)
		}
	}
	fin := columnize.SimpleFormat(output)
	fmt.Println(fin)
}

func (m *MusicBrainz) SearchRelease(q *ReleaseQuerry) error {
	querry := fmt.Sprintf("release:%s, artist:%s", q.album, q.artist)

	res, err := m.client.SearchRelease(querry, -1, -1)
	if err != nil {
		return err
	}
	m.ReleaseSearchResponses = append(m.ReleaseSearchResponses, res)
	return nil
}

func (m *MusicBrainz) GetQuerry() {
	reader := bufio.NewReader(os.Stdin)

	// get album name
	fmt.Println("Please enter the album name")
	album, err := reader.ReadString('\n')
	if err != nil {
		fmt.Print(err)
	}
	album = strings.Replace(album, "\n", "", -1)

	fmt.Println("Please enter the artist name")
	artist, err := reader.ReadString('\n')
	if err != nil {
		fmt.Print(err)
	}
	artist = strings.Replace(artist, "\n", "", -1)

	fmt.Printf("Album: %s Artist: %s \n", album, artist)

	querryData := ReleaseQuerry{
		artist: artist,
		album:  album,
	}
	m.ReleaseQuerrys = append(m.ReleaseQuerrys, querryData)
}
