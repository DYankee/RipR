package musicbrainz

type Release struct {
	ID     string
	Title  string
	Artist string
	Date   string
	Tracks []Track
}

type Track struct {
	Position int
	Title    string
	Length   int //milliseconds
}
