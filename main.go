package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/michiwend/gomusicbrainz"
)

type AlbumSearchQuerry struct {
	album  string
	artist string
}

func (a *AlbumSearchQuerry) getAlbumName() {
	fmt.Scanln(&a.album)
}

func (a *AlbumSearchQuerry) getArtistName() {
	fmt.Scanln(&a.artist)
}

func (a AlbumSearchQuerry) newSearch() {
	fmt.Println("Please enter the album name")
	a.getAlbumName()
	fmt.Println("Please enter the artist name")
	a.getArtistName()
}

type Musicbrainz struct {
	client              *gomusicbrainz.WS2Client
	ReleaseSearchQuerry AlbumSearchQuerry
	RleaseSearchRes     *gomusicbrainz.ReleaseSearchResponse
}

func (m *Musicbrainz) newInst() {
	var err error
	m.client, err = gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"RRipper",
		"0.0.1-beta",
		"http://github.com/Dyankee/RRipper")
	if err != nil {
		fmt.Println(err)
	}
}

func (m *Musicbrainz) ReleaseSearch() {
	var err error
	search := fmt.Sprintf(`release:%s, artist:%s`, m.ReleaseSearchQuerry.album, m.ReleaseSearchQuerry.artist)
	m.RleaseSearchRes, err = m.client.SearchRelease(search, -1, -1)
	if err != nil {
		fmt.Println(err)
	}
}

func DisplayReleaseRes(res *gomusicbrainz.ReleaseSearchResponse) {
	var id int
	for _, release := range res.Releases {
		id += 1
		if release.Mediums[0].Format == `12" Vinyl` {
			fmt.Printf("%-3d %s Name: %-15s  Artist: %-15s  Format %-15s  Release Date: %d Country: %s\n",
				id,
				release.ID,
				release.Title,
				release.ArtistCredit.NameCredits[0].Artist.Name,
				release.Mediums[0].Format,
				release.Date.Year(),
				release.CountryCode,
			)
		}
	}
}

func getQuerryterms(s *bufio.Reader) (album string, artist string) {
	fmt.Println("Please enter the album name")
	album, err := s.ReadString('\n')
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("Please enter the artist name")
	artist, err = s.ReadString('\n')
	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("Album: %s Artist %s \n", album, artist)
	return album, artist
}

func buildReleaseQuerry(album string, artist string) string {
	querry := fmt.Sprintf(`release:%s, artist:%s`, album, artist)
	return querry
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	client, err := gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"RRipper",
		"0.0.1-beta",
		"http://github.com/Dyankee/RRipper")
	if err != nil {
		fmt.Println(err)
	}

	res, err := client.SearchRelease(buildReleaseQuerry(getQuerryterms(reader)), -1, -1)
	if err != nil {
		fmt.Println(err)
	}

	DisplayReleaseRes(res)

	//choice := resp.Releases[8-1].ID
	//resp2, err := Musicbrainz.client.LookupRelease(choice, "media+recordings")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Printf("Track %d: %s Length: %d ID: %s\n",
	//	resp2.Mediums[0].Tracks[0].Position,
	//	resp2.Mediums[0].Tracks[0].Recording.Title,
	//	resp2.Mediums[0].Tracks[0].Length,
	//	resp2.ID,
	//)

	//audacity code
	//audacity := a.Audacity{}
	//audacity.Connect()
	//res := audacity.Do_command("StoreCursorPosition:")
	//fmt.Println(res)
}
