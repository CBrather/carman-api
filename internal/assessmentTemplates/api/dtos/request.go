package dtos

type Request struct {
	Label  string  `json:"label"`
	Name   string  `json:"name"`
	Scales []Scale `json:"scales"`
}
