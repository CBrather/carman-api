package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/CBrather/carman-api/internal/radarChartDesigns/api/dtos"
	"github.com/CBrather/carman-api/internal/radarChartDesigns/repository"
)

func GetByID(repo *repository.Design) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := chi.URLParam(req, "id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			http.Error(w, "Invalid id provided", http.StatusBadRequest)
			return
		}

		id64 := int64(id)
		album, err := repo.GetByID(req.Context(), id64)

		if err != nil {
			http.Error(w, "No album with that id was found", http.StatusNotFound)
			return
		}

		body, err := json.Marshal(album)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to serialize the album with id %s", album.ID), zap.Error(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}

func Create(repo *repository.Design) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var newDesign dtos.NewDesign

		rawRequestBody, err := io.ReadAll(req.Body)
		if err != nil {
			zap.L().Warn("Unable to read bytes of the request body", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(rawRequestBody, &newDesign)
		if err != nil {
			zap.L().Info("Unable to deserialize request body to album", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		newDesignModel := repository.DesignModel{Name: newDesign.Name}

		addedAlbum, err := repo.Create(req.Context(), newDesignModel)
		if err != nil {
			zap.L().Error("Failed to save new album", zap.Error(err))
			zap.L().Debug("Failed to save new album", zap.Any("struct", newDesign))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		serializedAlbum, err := json.Marshal(addedAlbum)
		if err != nil {
			zap.L().Error("Failed serializing new album after successful save", zap.Error(err))
			zap.L().Debug("Failed serializing new album after successful save", zap.Any("struct", addedAlbum))

			http.Error(w, "Internal Server Error occurred after the album was successfully saved", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(serializedAlbum)
	}
}
