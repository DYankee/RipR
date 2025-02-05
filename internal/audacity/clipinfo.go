package Audacity

type ClipInfo struct {
	Track int     `json:"track"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Color int     `json:"color"`
	Name  string  `json:"name"`
}

func (ci *ClipInfo) GetClipLength() float64 {
	return ci.End - ci.Start
}
