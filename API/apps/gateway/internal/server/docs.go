package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mvrilo/go-redoc"
)

func (s *Server) redocRoutes() http.Handler {
	r := chi.NewRouter()

	doc := redoc.Redoc{
		Title:       "BirthdayWish API",
		Description: "Gateway API for BirthdayWish",
		SpecFile:    "internal/api/openapi.yaml",
		SpecPath:    "/docs/openapi.yaml",
		DocsPath:    "/docs",
	}

	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/api/openapi.yaml")
	})

	r.Get("/", doc.Handler())

	return r
}
