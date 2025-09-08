// apps/gateway/internal/handlers/api_impl.go
package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/client"
	api "github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/gen"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/handlers/auth_handler"
	"github.com/go-chi/jwtauth/v5"
)

type APIImplementation struct {
	authHandler *auth_handler.AuthHandler
	authClient  *client.AuthClient
}

type HealthHandler struct{}

func NewAPIImplementation(authClient *client.AuthClient, log *slog.Logger, tokenAuth *jwtauth.JWTAuth) *APIImplementation {
	return &APIImplementation{
		authHandler: auth_handler.NewAuthHandler(authClient, log),
	}
}

func (a *APIImplementation) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	a.authHandler.SignIn(w, r)
}

func (a *APIImplementation) PostAuthSignup(w http.ResponseWriter, r *http.Request) {
	a.authHandler.SignUp(w, r)
}

func (a *APIImplementation) GetTestAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Authentication successful",
	})
}

var _ api.ServerInterface = (*APIImplementation)(nil)
