package api

import (
	"github.com/CBrather/carman-api/internal/radarChartDesigns/api/handlers"
	"github.com/CBrather/carman-api/internal/radarChartDesigns/repository"
	"github.com/CBrather/carman-api/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func Router(rootRouter *chi.Mux, dbClient *mongo.Client, jwtValidatorConfig middleware.JWTValidatorConfig) {
	router := chi.NewRouter()
	router.Use(middleware.EnsureValidToken(jwtValidatorConfig))

	repo, err := repository.New(dbClient)
	if err != nil {
		zap.L().Fatal("Failed creating RadarChartDesign Repository.", zap.Error(err))
	}

	router.With(middleware.RequirePermission("chartdesign:create")).Post("/", handlers.Create(repo))
	router.With(middleware.RequirePermission("chartdesign:list")).Get("/", handlers.List(repo))
	router.With(middleware.RequirePermission("chartdesign:read")).Get("/{id}", handlers.GetByID(repo))
	router.With(middleware.RequirePermission("chartdesign:update")).Put("/{id}", handlers.UpdateByID(repo))

	rootRouter.Mount("/charts/designs/radar", router)
}
