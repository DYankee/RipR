package main

import (
	"fmt"

	"github.com/michiwend/gomusicbrainz"
)

func main() {

	client, err := gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"A GoMusicBrainz example",
		"0.0.1-beta",
		"http://github.com/michiwend/gomusicbrainz")
	if err != nil {
		fmt.Println(err)
	}
	album := "let it be"
	artist := "the beatles"
	search := fmt.Sprintf(`release:%s, artist:%s`, album, artist)
	resp, err := client.SearchRelease(search, -1, -1)
	if err != nil {
		fmt.Println(err)
	}
	var id int
	for _, release := range resp.Releases {
		if release.Mediums[0].Format == `12" Vinyl` {
			id += 1
			fmt.Printf("%s Name: %-20s  Artist: %-20s  Format %-15s  Release Date: %d Country: %s\n",
				release.ID,
				release.Title,
				release.ArtistCredit.NameCredits[0].Artist.Name,
				release.Mediums[0].Format,
				release.Date.Year(),
				release.CountryCode,
			)
		}
	}
	choice := resp.Releases[0].ID

	resp2, err := client.LookupRelease(choice)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s %s",
		resp2.Title,
		resp2.ArtistCredit.NameCredits[0].Artist.Name,
	)

	//audacity code
	//audacity := a.Audacity{}
	//audacity.Connect()
	//res := audacity.Do_command("StoreCursorPosition:")
	//fmt.Println(res)
}
