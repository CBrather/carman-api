package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type DesignModel struct {
	ID     string  `bson:"_id"`
	Name  string  `bson:"name"`
	CircularEdges edgeDesign  `bson:"circularEdges"`
	OuterEdge  edgeDesign `bson:"outerEdge"`
	RadialEdges edgeDesign `bson:"radialEdges"`
}

type edgeDesign struct {
	Color string `bson:"color"`
	Style string `bson:"style"`
	Thickness int `bson:"thickness"`
}

type Design struct {
	DBClient *mongo.Client
}

func New(dbClient *mongo.Client) (*Design, error) {
	return &Design{DBClient: dbClient}, nil
}

func (r *Design) Create(ctx context.Context, newDesign DesignModel) (DesignModel, error) {
	return DesignModel{}, nil
}

func (r *Design) GetByID(ctx context.Context, id int64) (DesignModel, error) {
	return DesignModel{}, nil
}

func (r *Design) List(ctx context.Context) ([]DesignModel, error) {
	return []DesignModel{{}}, nil
}

