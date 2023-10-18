package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/CBrather/carman-api/internal/radarChartDesigns/repository"
)

func GetByID(repo *repository.Design) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		design, err := repo.GetByID(req.Context(), id)
		if err != nil {
			http.Error(w, "No design with that id was found", http.StatusNotFound)
			return
		}

		body, err := json.Marshal(design)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to serialize the design with id %s", design.ID), zap.Error(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}
