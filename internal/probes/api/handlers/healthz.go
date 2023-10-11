package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupProbeRoutes(router chi.Router) {
	router.Get("/healthz", getHealth)
}

func getHealth(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write([]byte{})
}
