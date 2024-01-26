package api

import (
	"github.com/CBrather/carman-api/internal/radarChartDesigns/api/handlers"
	"github.com/CBrather/carman-api/internal/radarChartDesigns/repository"
	"github.com/CBrather/carman-api/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type RouterConfig struct {
	DBClient *mongo.Client
	Auth     middleware.JWTValidatorConfig
}

func Router(rootRouter *chi.Mux, config RouterConfig) {
	router := chi.NewRouter()
	router.Use(middleware.EnsureValidToken(config.Auth))

	repo, err := repository.New(config.DBClient)
	if err != nil {
		zap.L().Fatal("Failed creating RadarChartDesign Repository.", zap.Error(err))
	}

	router.With(middleware.RequirePermission("chartdesign:create")).Post("/", handlers.Create(repo))
	router.With(middleware.RequirePermission("chartdesign:list")).Get("/", handlers.List(repo))
	router.With(middleware.RequirePermission("chartdesign:read")).Get("/{id}", handlers.GetByID(repo))

	rootRouter.Mount("/charts/designs/radar", router)
}
