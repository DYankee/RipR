package main

import (
	Internal "github.com/DYankee/RRipper/internal"
)

func main() {
	a := Internal.Audacity{}
	a.Open("Thriller.aup3")

	a.Connect()
	// c := exec.Command("audacity", "Thriller.aup3")
	// c.Run()
	//	mb := Internal.MusicBrainz{}
	//	mb.Init()
	//	mb.GetQuerry()
	//	mb.SearchRelease(&mb.ReleaseQuerrys[0])

	//	fmt.Println(mb.ReleaseSearchResponses[0].Releases[0].Title)

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

	//audacity.Connect()
	//res := audacity.Do_command("StoreCursorPosition:")
	//fmt.Println(res)
}
