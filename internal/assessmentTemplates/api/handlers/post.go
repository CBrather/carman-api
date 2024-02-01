package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/CBrather/carman-api/internal/assessmentTemplates/api/dtos"
	"github.com/CBrather/carman-api/internal/assessmentTemplates/repository"
)

func Create(repo *repository.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rawRequestBody, err := io.ReadAll(req.Body)
		if err != nil {
			zap.L().Warn("Unable to read bytes of the request body", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var newTemplate dtos.Request
		if err = json.Unmarshal(rawRequestBody, &newTemplate); err != nil {
			zap.L().Info("Unable to deserialize request body to template", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		newScaleModels := make([]repository.ScaleModel, 0, len(newTemplate.Scales))
		for _, scale := range newTemplate.Scales {
			newScaleModels = append(newScaleModels, scale.ToModel())
		}

		newTemplateModel := repository.TemplateModel{
			Label: newTemplate.Label,
			Name:  newTemplate.Name,
		}

		createdTemplate, err := repo.Create(req.Context(), newTemplateModel, newScaleModels)
		if err != nil {
			zap.L().Error("Failed to save new template", zap.Error(err))
			zap.L().Debug("Failed to save new template", zap.Any("struct", newTemplate))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		createdScales, err := repo.ScalesRepo.ListByID(req.Context(), createdTemplate.Scales)
		if err != nil {
			zap.L().Error("Failed to get scales for new template after successful save", zap.Error(err))
			zap.L().Debug("Failed to get scales for new template after successful save", zap.Any("struct", createdTemplate))

			http.Error(w, "Internal Server Error occurred after the template was successfully saved", http.StatusInternalServerError)
			return
		}

		responseDTO := dtos.Response{
			ID:     createdTemplate.IDString(),
			Label:  createdTemplate.Label,
			Name:   createdTemplate.Name,
			Scales: dtos.ScalesFromModel(createdScales),
		}

		responseBody, err := json.Marshal(responseDTO)
		if err != nil {
			zap.L().Error("Failed serializing new template after successful save", zap.Error(err))
			zap.L().Debug("Failed serializing new template after successful save", zap.Any("struct", createdTemplate))

			http.Error(w, "Internal Server Error occurred after the template was successfully saved", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(responseBody)
	}
}
