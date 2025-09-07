package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/config"
	api "github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/gen"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/handlers/auth_handler"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/storage"
	"github.com/go-chi/jwtauth/v5"
)

type APIImplementation struct {
	authHandler   *auth_handler.AuthHandler
	healthHandler *HealthHandler
}

type HealthHandler struct{}

func (a *APIImplementation) GetHealth(w http.ResponseWriter, r *http.Request) {
	a.healthHandler.GetHealth(w, r)
}

func (a *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	status := "healthy"
	timestamp := time.Now()
	response := api.HealthResponse{
		Status:    &status,
		Timestamp: &timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func NewAPIImplementation(cfg *config.Config, storage *storage.Storage, log *slog.Logger, tokenAuth *jwtauth.JWTAuth) *APIImplementation {
	return &APIImplementation{
		authHandler: auth_handler.NewAuthHandler(cfg, storage, log, tokenAuth),
	}
}

func (a *APIImplementation) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	a.authHandler.SignIn(w, r)
}

func (a *APIImplementation) PostAuthSignup(w http.ResponseWriter, r *http.Request) {
	a.authHandler.SignUp(w, r)
}

func (a *APIImplementation) GetTestAuth(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"message": "Authentication successful",
		"claims":  claims,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

var _ api.ServerInterface = (*APIImplementation)(nil)
