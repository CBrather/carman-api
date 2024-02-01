package dtos

type Response struct {
	ID     string  `json:"id"`
	Label  string  `json:"label"`
	Name   string  `json:"name"`
	Scales []Scale `json:"scales"`
}

type ResponseList struct {
	Items []Response `json:"items"`
}
