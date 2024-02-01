package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/CBrather/carman-api/internal/radarChartDesigns/api/dtos"
	"github.com/CBrather/carman-api/internal/radarChartDesigns/repository"
)

func Create(repo *repository.Design) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rawRequestBody, err := io.ReadAll(req.Body)
		if err != nil {
			zap.L().Warn("Unable to read bytes of the request body", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var newDesign dtos.ChartDesignRequest
		if err = json.Unmarshal(rawRequestBody, &newDesign); err != nil {
			zap.L().Info("Unable to deserialize request body to design", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		newDesignModel := repository.DesignModel{
			Name:          newDesign.Name,
			CircularEdges: repository.EdgeDesign(newDesign.CircularEdges),
			OuterEdge:     repository.EdgeDesign(newDesign.OuterEdge),
			RadialEdges:   repository.EdgeDesign(newDesign.RadialEdges),
			StartingAngle: newDesign.StartingAngle,
		}

		createdDesign, err := repo.Create(req.Context(), newDesignModel)
		if err != nil {
			zap.L().Error("Failed to save new design", zap.Error(err))
			zap.L().Debug("Failed to save new design", zap.Any("struct", newDesign))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		responseDTO := dtos.ResponseSingle{
			ID:            createdDesign.IDString(),
			Name:          createdDesign.Name,
			CircularEdges: dtos.EdgeDesign(createdDesign.CircularEdges),
			OuterEdge:     dtos.EdgeDesign(createdDesign.OuterEdge),
			RadialEdges:   dtos.EdgeDesign(createdDesign.RadialEdges),
			StartingAngle: createdDesign.StartingAngle,
		}

		responseBody, err := json.Marshal(responseDTO)
		if err != nil {
			zap.L().Error("Failed serializing new design after successful save", zap.Error(err))
			zap.L().Debug("Failed serializing new design after successful save", zap.Any("struct", createdDesign))

			http.Error(w, "Internal Server Error occurred after the design was successfully saved", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(responseBody)
	}
}
