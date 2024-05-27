package main

import (
	Internal "github.com/DYankee/RRipper/internal"
)

func main() {
	a := Internal.Audacity{}
	a.Open("testdata/Thriller.aup3")
	a.Connect()

	a.SelectRegion(0, 10)
	a.ExportAudio("/home/z-geary/code/rripper/testdata", "test.flac")

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
