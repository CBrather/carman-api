package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/CBrather/carman-api/internal/app_errors"
	"github.com/CBrather/carman-api/internal/assessmentTemplates/api/dtos"
	"github.com/CBrather/carman-api/internal/assessmentTemplates/repository"
)

func GetByID(repo *repository.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		template, err := repo.GetByID(req.Context(), id)
		if err != nil {
			if errors.Is(err, app_errors.ErrNotFound{}) {
				zap.L().Warn(fmt.Sprintf("Failed getting template, as none with id %s was found", id))
				http.Error(w, "No template with id was found", http.StatusNotFound)
			} else {
				zap.L().Error(fmt.Sprintf("Failed getting template with id %s", id), zap.Error(err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			return
		}

		scales, err := repo.ScalesRepo.ListByID(req.Context(), template.Scales)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed getting scales for template with id %s", id), zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		scalesDTO := dtos.ScalesFromModel(scales)

		responseDTO := dtos.Response{
			ID:     template.IDString(),
			Label:  template.Label,
			Name:   template.Name,
			Scales: scalesDTO,
		}

		body, err := json.Marshal(responseDTO)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to serialize the template with id %s", template.ID), zap.Error(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}

func List(repo *repository.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		templates, err := repo.List(req.Context())
		if err != nil {
			http.Error(w, "failed retrieving templates", http.StatusInternalServerError)
			return
		}

		var responseDTO dtos.ResponseList
		for _, template := range templates {
			scales, err := repo.ScalesRepo.ListByID(req.Context(), template.Scales)
			if err != nil {
				zap.L().Error(fmt.Sprintf("Failed getting scales for template with id %s", template.ID.Hex()), zap.Error(err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			scalesDTO := dtos.ScalesFromModel(scales)
			templateDTO := dtos.Response{
				ID:     template.IDString(),
				Name:   template.Name,
				Label:  template.Label,
				Scales: scalesDTO,
			}

			responseDTO.Items = append(responseDTO.Items, templateDTO)
		}

		body, err := json.Marshal(responseDTO)
		if err != nil {
			zap.L().Error("Failed to serialize the templates", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}
