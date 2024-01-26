package dtos

type ChartDesignRequest struct {
	Name          string     `json:"name"`
	CircularEdges EdgeDesign `json:"circularEdges"`
	OuterEdge     EdgeDesign `json:"outerEdge"`
	RadialEdges   EdgeDesign `json:"radialEdges"`
	StartingAngle int        `json:"startingAngle"`
}
