package main

import (
	a "github.com/DYankee/RRipper/internal"
)

func main() {

	mb := a.MusicBrainz{}
	mb.Init()

	mb.GetQuerry()
	mb.SearchRelease(&mb.ReleaseQuerrys[0])
	mb.DisplayReleaseRes(mb.ReleaseSearchResponses[0])

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
