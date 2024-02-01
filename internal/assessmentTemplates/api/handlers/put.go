package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/CBrather/carman-api/internal/app_errors"
	"github.com/CBrather/carman-api/internal/assessmentTemplates/api/dtos"
	"github.com/CBrather/carman-api/internal/assessmentTemplates/repository"
	"github.com/go-chi/chi/v5"
)

func UpdateByID(repo *repository.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if id == "" {
			zap.L().Warn("UpdateByID handler called without id param")

			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

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

		updatedTemplate, err := repo.UpdateByID(req.Context(), id, newTemplateModel, newScaleModels)
		if err != nil {
			if errors.Is(err, app_errors.ErrNotFound{}) {
				zap.L().Warn(fmt.Sprintf("Failed updating template, as none with id %s was found", id))
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else {
				zap.L().Error(fmt.Sprintf("Failed to save updated template with id %s", id), zap.Error(err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			return
		}

		updatedScales, err := repo.ScalesRepo.ListByID(req.Context(), updatedTemplate.Scales)
		if err != nil {
			zap.L().Error("Failed to get scales for updated template after successful save", zap.Error(err))
			zap.L().Debug("Failed to get scales for updated template after successful save", zap.Any("struct", updatedTemplate))
			http.Error(w, "Internal Server Error occurred after the template was successfully updated", http.StatusInternalServerError)
			return
		}

		responseDTO := dtos.Response{
			ID:     updatedTemplate.IDString(),
			Label:  updatedTemplate.Label,
			Name:   updatedTemplate.Name,
			Scales: dtos.ScalesFromModel(updatedScales),
		}

		responseBody, err := json.Marshal(responseDTO)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed serializing updated template with id %s after successful save", id), zap.Error(err))

			http.Error(w, "Internal Server Error occurred after the template was successfully saved", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(responseBody)
	}
}
