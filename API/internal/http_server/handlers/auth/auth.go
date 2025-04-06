package authhandler

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/storage"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type AuthHandler struct {
	cfg     *config.Config
	storage *storage.Storage
	log     *slog.Logger
}

type IAuthHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	// SignIn(w http.ResponseWriter, r *http.Request)
	// Refresh(w http.ResponseWriter, r *http.Request)
}

func New(cfg *config.Config, storage *storage.Storage, log *slog.Logger) *AuthHandler {
	return &AuthHandler{cfg, storage, log}
}

func validateCredentials(user Credentials) error {
	var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)
	var passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()-_+=]{8,20}$`)

	if !usernameRegex.MatchString(user.Username) || !passwordRegex.MatchString(user.Password) {
		return errors.New("invalid username or password")
	}

	return nil
}
