// models/messages.go
package models

import (
	"github.com/michiwend/gomusicbrainz"
)

// SearchSubmittedMsg is sent when the user clicks submit
type SearchSubmittedMsg struct {
	Artist string
	Album  string
}

type SearchResultMsg struct {
	Response gomusicbrainz.ReleaseSearchResponse
}

type SearchErrMsg struct {
	Err error
}
