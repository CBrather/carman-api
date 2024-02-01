package api

import (
	"github.com/CBrather/carman-api/internal/assessmentTemplates/api/handlers"
	"github.com/CBrather/carman-api/internal/assessmentTemplates/repository"
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

	router.Post("/", handlers.Create(repo))
	router.Get("/", handlers.List(repo))
	router.Get("/{id}", handlers.GetByID(repo))
	router.Put("/{id}", handlers.UpdateByID(repo))

	rootRouter.Mount("/assessmentTemplates", router)
}
