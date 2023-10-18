package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/CBrather/carman-api/internal/radarChartDesigns/api/dtos"
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

		responseDTO := dtos.ResponseSingle{
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
			designDTO := dtos.ResponseSingle{
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
