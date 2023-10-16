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
	StartingAngle int                `bson:"int"`
}

type EdgeDesign struct {
	Color     string `bson:"color"`
	Style     string `bson:"style"`
	Thickness int    `bson:"thickness"`
}

type Design struct {
	DBCollection *mongo.Collection
}

func New(dbClient *mongo.Client) (*Design, error) {
	return &Design{DBCollection: dbClient.Database("carman").Collection("radar_chart_designs")}, nil
}

func (r *Design) Create(ctx context.Context, newDesign DesignModel) (*DesignModel, error) {
	result, err := r.DBCollection.InsertOne(context.TODO(), newDesign)
	if err != nil {
		return nil, err
	}

	var findResult DesignModel
	err = r.DBCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&findResult)
	if err != nil {
		return nil, err
	}

	return &findResult, nil
}

func (r *Design) GetByID(ctx context.Context, id int64) (DesignModel, error) {
	return DesignModel{}, nil
}

func (r *Design) List(ctx context.Context) ([]DesignModel, error) {
	return []DesignModel{{}}, nil
}
