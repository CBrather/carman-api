package dtos

import (
	"fmt"

	"github.com/CBrather/carman-api/internal/assessmentTemplates/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type Scale struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Name  string `json:"name"`
	Steps []step `json:"steps"`
}

type step struct {
	Label string `json:"label"`
	Value int32  `json:"value"`
}

func (s step) ToModel() repository.ScaleStep {
	return repository.ScaleStep{
		Label: s.Label,
		Value: s.Value,
	}
}

func (s Scale) ToModel() repository.ScaleModel {
	steps := make([]repository.ScaleStep, 0, len(s.Steps))
	for _, step := range s.Steps {
		steps = append(steps, step.ToModel())
	}

	model := repository.ScaleModel{
		Label: s.Label,
		Name:  s.Name,
		Steps: steps,
	}

	if objectID, err := primitive.ObjectIDFromHex(s.ID); err == nil {
		model.ID = objectID
	} else {
		zap.L().Debug(fmt.Sprintf("Could not convert scale ID %s to objectID, continuing with nil ID", s.ID), zap.Error(err))
	}

	return model
}

func ScalesFromModel(models []repository.ScaleModel) []Scale {
	scales := make([]Scale, 0, len(models))
	for _, model := range models {
		scales = append(scales, ScaleFromModel(model))
	}

	return scales
}

func ScaleFromModel(model repository.ScaleModel) Scale {
	steps := make([]step, 0, len(model.Steps))
	for _, step := range model.Steps {
		steps = append(steps, stepFromModel(step))
	}

	return Scale{
		ID:    model.IDString(),
		Label: model.Label,
		Name:  model.Name,
		Steps: steps,
	}
}

func stepFromModel(model repository.ScaleStep) step {
	return step{
		Label: model.Label,
		Value: model.Value,
	}
}
