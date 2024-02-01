package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	assessmentTemplates "github.com/CBrather/carman-api/internal/assessmentTemplates/api"
	probes "github.com/CBrather/carman-api/internal/probes/api/handlers"
	radarChartDesigns "github.com/CBrather/carman-api/internal/radarChartDesigns/api"
	"github.com/CBrather/carman-api/pkg/database"
	"github.com/CBrather/carman-api/pkg/middleware"
)

func SetupHttpRoutes(env EnvConfig) {
	dbClient, err := database.GetDBClient(env.DBConnectionString)
	if err != nil {
		zap.L().Fatal("failed connecting to the database", zap.Error(err))
	}

	zap.L().Info("database connection is established")

	router := chi.NewRouter()

	router.Use(
		chiMiddleware.Recoverer,
	)

	probes.SetupProbeRoutes(router)

	jwtValidatorConfig := middleware.JWTValidatorConfig{
		Audience: env.Auth.Audience,
		Domain:   env.Auth.Domain,
	}
	assessmentTemplates.Router(router, dbClient, jwtValidatorConfig)
	radarChartDesigns.Router(router, dbClient, jwtValidatorConfig)

	zap.L().Info("Server listening on :8080")

	err = http.ListenAndServe("0.0.0.0:8080", router)
	if err != nil {
		zap.L().Fatal("Server threw an error")
	}
}
