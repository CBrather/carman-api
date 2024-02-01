package dtos

type Response struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	CircularEdges EdgeDesign `json:"circularEdges"`
	OuterEdge     EdgeDesign `json:"outerEdge"`
	RadialEdges   EdgeDesign `json:"radialEdges"`
	StartingAngle int        `json:"startingAngle"`
}

type ResponseList struct {
	Items []Response `json:"items"`
}
