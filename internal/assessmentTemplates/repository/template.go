package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/CBrather/carman-api/internal/app_errors"
)

type TemplateModel struct {
	ID     primitive.ObjectID   `bson:"_id,omitempty"`
	Label  string               `bson:"label"`
	Name   string               `bson:"name"`
	Scales []primitive.ObjectID `bson:"scales"`
}

func (m TemplateModel) IDString() string {
	return m.ID.Hex()
}

type Template struct {
	dbCollection *mongo.Collection
	ScalesRepo   *Scale
}

func New(dbClient *mongo.Client) (*Template, error) {
	scalesRepo, err := newScaleRepository(dbClient)
	if err != nil {
		return nil, err
	}

	return &Template{
		dbCollection: dbClient.Database("carman").Collection("assessment_templates"),
		ScalesRepo:   scalesRepo}, nil
}

func (r *Template) Create(ctx context.Context, newTemplate TemplateModel, scales []ScaleModel) (*TemplateModel, error) {
	for _, scale := range scales {
		scaleResult, err := r.ScalesRepo.Create(ctx, scale)
		if err != nil {
			return nil, err
		}

		newTemplate.Scales = append(newTemplate.Scales, scaleResult.ID)
	}

	result, err := r.dbCollection.InsertOne(context.TODO(), newTemplate)
	if err != nil {
		return nil, err
	}

	var findResult TemplateModel
	err = r.dbCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&findResult)
	if err != nil {
		return nil, err
	}

	return &findResult, nil
}

func (r *Template) GetByID(ctx context.Context, id string) (*TemplateModel, error) {
	var foundTemplate TemplateModel

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.dbCollection.FindOne(ctx, bson.D{{Key: "_id", Value: objectID}}).Decode(&foundTemplate)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, app_errors.NewErrNotFound(err)
	}

	return &foundTemplate, err
}

func (r *Template) List(ctx context.Context) ([]TemplateModel, error) {
	cursor, err := r.dbCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var foundTemplates []TemplateModel
	err = cursor.All(ctx, &foundTemplates)

	return foundTemplates, err
}

func (r *Template) UpdateByID(ctx context.Context, id string, newTemplate TemplateModel, scales []ScaleModel) (*TemplateModel, error) {
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	newScales := make([]primitive.ObjectID, 0, len(scales))

	for _, scale := range scales {
		if scale.ID.IsZero() {
			scaleResult, err := r.ScalesRepo.Create(ctx, scale)
			if err != nil {
				return nil, err
			}

			newScales = append(newScales, scaleResult.ID)
		} else {
			scaleResult, err := r.ScalesRepo.UpdateByID(ctx, scale.IDString(), scale)
			if err != nil {
				return nil, err
			}

			newScales = append(newScales, scaleResult.ID)
		}
	}

	newTemplate.Scales = newScales

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.D{{Key: "$set", Value: newTemplate}}
	_, err = r.dbCollection.UpdateByID(context.TODO(), objectID, update)
	if err != nil {
		return nil, err
	}

	var findResult TemplateModel
	err = r.dbCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectID}}).Decode(&findResult)
	if err != nil {
		return nil, err
	}

	return &findResult, nil
}
