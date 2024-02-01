package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/CBrather/carman-api/internal/app_errors"
)

type ScaleStep struct {
	Label string `bson:"label"`
	Value int32  `bson:"value"`
}

type ScaleModel struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Label string             `bson:"label"`
	Name  string             `bson:"name"`
	Steps ScaleStep          `bson:"steps"`
}

func (m ScaleModel) IDString() string {
	return m.ID.Hex()
}

type Scale struct {
	DBCollection *mongo.Collection
}

func newScaleRepository(dbClient *mongo.Client) (*Scale, error) {
	return &Scale{DBCollection: dbClient.Database("carman").Collection("scales")}, nil
}

func (r *Scale) Create(ctx context.Context, newScale ScaleModel) (*ScaleModel, error) {
	result, err := r.DBCollection.InsertOne(context.TODO(), newScale)
	if err != nil {
		return nil, err
	}

	var findResult ScaleModel
	err = r.DBCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&findResult)
	if err != nil {
		return nil, err
	}

	return &findResult, nil
}

func (r *Scale) GetByID(ctx context.Context, id string) (*ScaleModel, error) {
	var scale ScaleModel

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.DBCollection.FindOne(ctx, bson.D{{Key: "_id", Value: objectID}}).Decode(&scale)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, app_errors.NewErrNotFound(err)
	}

	return &scale, err
}

func (r *Scale) List(ctx context.Context) ([]ScaleModel, error) {
	cursor, err := r.DBCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var scales []ScaleModel
	err = cursor.All(ctx, &scales)

	return scales, err
}

func (r *Scale) UpdateByID(ctx context.Context, id string, newScale ScaleModel) (*ScaleModel, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.D{{Key: "$set", Value: newScale}}
	updateResult, err := r.DBCollection.UpdateByID(context.TODO(), objectID, update)
	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, app_errors.NewErrNotFound(nil)
	}

	var findResult ScaleModel
	err = r.DBCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectID}}).Decode(&findResult)
	if err != nil {
		return nil, err
	}

	return &findResult, nil
}
