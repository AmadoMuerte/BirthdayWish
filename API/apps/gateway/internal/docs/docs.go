package docs

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/mvrilo/go-redoc"
)

func RedocRoutes() http.Handler {
	r := chi.NewRouter()

	specFile := filepath.Join("apps", "gateway", "internal", "docs", "openapi.yaml")

	doc := redoc.Redoc{
		Title:    "BirthdayWish API",
		SpecFile: specFile,
		SpecPath: "/docs/openapi.yaml",
		DocsPath: "",
	}

	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, specFile)
	})

	r.Get("/*", doc.Handler())

	return r
}
