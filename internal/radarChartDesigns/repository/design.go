package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DesignModel struct {
	ID            primitive.ObjectID `bson:"_id, omitempty"`
	Name          string             `bson:"name"`
	CircularEdges EdgeDesign         `bson:"circularEdges"`
	OuterEdge     EdgeDesign         `bson:"outerEdge"`
	RadialEdges   EdgeDesign         `bson:"radialEdges"`
}

type EdgeDesign struct {
	Color     string `bson:"color"`
	Style     string `bson:"style"`
	Thickness int    `bson:"thickness"`
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
