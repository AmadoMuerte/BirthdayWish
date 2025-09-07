package auth_handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/config"
	gen "github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/gen"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/jwt"
	models "github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/models"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/response"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/storage"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	cfg       *config.Config
	storage   *storage.Storage
	log       *slog.Logger
	tokenAuth *jwtauth.JWTAuth
}

func NewAuthHandler(cfg *config.Config, storage *storage.Storage, log *slog.Logger, tokenAuth *jwtauth.JWTAuth) *AuthHandler {
	return &AuthHandler{
		cfg:       cfg,
		storage:   storage,
		log:       log,
		tokenAuth: tokenAuth,
	}
}

func validateCredentials(validateType string, data string) error {
	var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)
	var passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()-_+=]{8,20}$`)
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if validateType == "username" {
		if !usernameRegex.MatchString(data) {
			return errors.New("invalid username")
		}
	}
	if validateType == "password" {
		if !passwordRegex.MatchString(data) {
			return errors.New("invalid password")
		}
	}
	if validateType == "email" {
		if !emailRegex.MatchString(data) {
			return errors.New("invalid email")
		}
	}

	return nil
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req gen.SignUpRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := validateCredentials("username", req.Username); err != nil {
		h.log.Error("invalid credentials", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := validateCredentials("password", req.Password); err != nil {
		h.log.Error("invalid password", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := validateCredentials("email", string(req.Email)); err != nil {
		h.log.Error("invalid email", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	exists, err := h.storage.UserExists(ctx, req.Username, string(req.Email))
	if err != nil {
		h.log.Error("database error checking user existence", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	if exists {
		h.log.Error("user with this username or email already exists")
		response.ErrorResponseJSON(w, r, http.StatusConflict, "User already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.log.Error("failed to hash password", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	user := models.User{
		Name:     req.Username,
		Password: string(hashedPassword),
		Email:    string(req.Email),
	}

	if err := h.storage.CreateUser(ctx, &user); err != nil {
		h.log.Error("failed to create user", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	resMessage := "User created successfully"
	responseData := gen.SignUpResponse{
		Message: &resMessage,
		UserId:  &user.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseData)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req gen.SignInRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("error decode credentials", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusUnauthorized, "Authentication failed")
		return
	}

	if err := validateCredentials("username", req.Username); err != nil {
		h.log.Error("invalid credentials", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusUnauthorized, "Authentication failed")
		return
	}

	if err := validateCredentials("password", req.Password); err != nil {
		h.log.Error("invalid password", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusUnauthorized, "Authentication failed")
		return
	}

	user, err := h.storage.GetUserByUsername(ctx, req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		h.log.Error("Error get user by username or password mismatch", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusUnauthorized, "Authentication failed")
		return
	}

	now := time.Now()
	jwtClaims := jwt.Claims{
		UserID: user.ID,
		Exp:    now.Add(24 * time.Hour).Unix(),
		Iat:    now.Unix(),
	}
	_, tokenString, err := h.tokenAuth.Encode(jwt.NewClaims(jwtClaims))
	if err != nil {
		h.log.Error("failed to generate token", "error", err, "userID", user.ID)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)

	exp := time.Now().Add(24 * time.Hour)
	responseData := gen.SignInResponse{
		Email:  &user.Email,
		Exp:    &exp,
		Name:   &user.Name,
		Token:  &tokenString,
		UserId: &user.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}
