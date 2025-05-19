package authhandler

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/storage"
	"github.com/go-chi/jwtauth/v5"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type AuthHandler struct {
	cfg       *config.Config
	storage   *storage.Storage
	log       *slog.Logger
	tokenAuth *jwtauth.JWTAuth
}

type IAuthHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	// Refresh(w http.ResponseWriter, r *http.Request)
}

func New(cfg *config.Config, storage *storage.Storage, log *slog.Logger) *AuthHandler {
	tokenAuth := jwtauth.New("HS256", []byte(cfg.App.SecretKey), nil)
	return &AuthHandler{cfg, storage, log, tokenAuth}
}

func validateCredentials(user Credentials) error {
	var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)
	var passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()-_+=]{8,20}$`)
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !usernameRegex.MatchString(user.Username) ||
		!passwordRegex.MatchString(user.Password) ||
		!emailRegex.MatchString(user.Email) {
		return errors.New("invalid username, password, or email")
	}

	return nil
}
