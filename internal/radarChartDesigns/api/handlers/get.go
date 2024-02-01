package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/CBrather/carman-api/internal/app_errors"
	"github.com/CBrather/carman-api/internal/radarChartDesigns/api/dtos"
	"github.com/CBrather/carman-api/internal/radarChartDesigns/repository"
)

func GetByID(repo *repository.Design) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		design, err := repo.GetByID(req.Context(), id)
		if err != nil {
			if errors.Is(err, app_errors.ErrNotFound{}) {
				zap.L().Warn(fmt.Sprintf("Failed getting design, as none with id %s was found", id))
				http.Error(w, "No design with id was found", http.StatusNotFound)
			} else {
				zap.L().Error(fmt.Sprintf("Failed getting design with id %s", id), zap.Error(err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			return
		}

		responseDTO := dtos.Response{
			ID:            design.IDString(),
			Name:          design.Name,
			CircularEdges: dtos.EdgeDesign(design.CircularEdges),
			OuterEdge:     dtos.EdgeDesign(design.OuterEdge),
			RadialEdges:   dtos.EdgeDesign(design.RadialEdges),
			StartingAngle: design.StartingAngle,
		}

		body, err := json.Marshal(responseDTO)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to serialize the design with id %s", design.ID), zap.Error(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}

func List(repo *repository.Design) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		designs, err := repo.List(req.Context())
		if err != nil {
			http.Error(w, "failed retrieving designs", http.StatusInternalServerError)
			return
		}

		var responseDTO dtos.ResponseList
		for _, design := range designs {
			designDTO := dtos.Response{
				ID:            design.IDString(),
				Name:          design.Name,
				CircularEdges: dtos.EdgeDesign(design.CircularEdges),
				OuterEdge:     dtos.EdgeDesign(design.OuterEdge),
				RadialEdges:   dtos.EdgeDesign(design.RadialEdges),
				StartingAngle: design.StartingAngle,
			}

			responseDTO.Items = append(responseDTO.Items, designDTO)
		}

		body, err := json.Marshal(responseDTO)
		if err != nil {
			zap.L().Error("Failed to serialize the designs", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}
