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
	var id = -1
	for _, release := range resp.Releases {
		id += 1
		if release.Mediums[0].Format == `12" Vinyl` {
			fmt.Printf("%d %s Name: %-20s  Artist: %-20s  Format %-15s  Release Date: %d Country: %s\n",
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
	choice := resp.Releases[8].ID

	resp2, err := client.LookupRelease(choice)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s %s \n",
		resp2.Title,
		resp2.Mediums,
	)

	//audacity code
	//audacity := a.Audacity{}
	//audacity.Connect()
	//res := audacity.Do_command("StoreCursorPosition:")
	//fmt.Println(res)
}
