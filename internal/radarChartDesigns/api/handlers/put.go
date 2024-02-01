package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/CBrather/carman-api/internal/app_errors"
	"github.com/CBrather/carman-api/internal/radarChartDesigns/api/dtos"
	"github.com/CBrather/carman-api/internal/radarChartDesigns/repository"
	"github.com/go-chi/chi/v5"
)

func UpdateByID(repo *repository.Design) http.HandlerFunc {
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

		updatedDesign, err := repo.UpdateByID(req.Context(), id, newDesignModel)
		if err != nil {
			if errors.Is(err, app_errors.ErrNotFound{}) {
				zap.L().Warn(fmt.Sprintf("Failed updating design, as none with id %s was found", id))
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else {
				zap.L().Error(fmt.Sprintf("Failed to save updated design with id %s", id), zap.Error(err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			return
		}

		responseDTO := dtos.Response{
			ID:            updatedDesign.IDString(),
			Name:          updatedDesign.Name,
			CircularEdges: dtos.EdgeDesign(updatedDesign.CircularEdges),
			OuterEdge:     dtos.EdgeDesign(updatedDesign.OuterEdge),
			RadialEdges:   dtos.EdgeDesign(updatedDesign.RadialEdges),
			StartingAngle: updatedDesign.StartingAngle,
		}

		responseBody, err := json.Marshal(responseDTO)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed serializing updated design with id %s after successful save", id), zap.Error(err))

			http.Error(w, "Internal Server Error occurred after the design was successfully saved", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(responseBody)
	}
}
