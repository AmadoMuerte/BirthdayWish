// apps/gateway/internal/handlers/auth_handler.go
package auth_handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/client"
	gen "github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/gen"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/response"
)

type AuthHandler struct {
	authClient *client.AuthClient
	log        *slog.Logger
}

func NewAuthHandler(authClient *client.AuthClient, log *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
		log:        log,
	}
}

func parseUserId(userId string) (int64, error) {
	return strconv.ParseInt(userId, 10, 64)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req gen.SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	grpcResp, err := h.authClient.SignUp(r.Context(), req.Username, string(req.Email), req.Password)
	if err != nil {
		h.log.Error("auth service error", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Registration failed")
		return
	}

	userId, err := parseUserId(grpcResp.UserId)
	if err != nil {
		h.log.Error("failed to parse user id", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Authentication failed")
		return
	}

	responseData := gen.SignUpResponse{
		Message: &grpcResp.Message,
		UserId:  &userId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseData)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req gen.SignInRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("error decode credentials", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusUnauthorized, "Authentication failed")
		return
	}

	grpcResp, err := h.authClient.SignIn(r.Context(), req.Username, req.Password)
	if err != nil {
		h.log.Error("auth service error", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusUnauthorized, "Authentication failed")
		return
	}

	exp := time.Unix(grpcResp.Exp, 0)
	w.Header().Set("Authorization", "Bearer "+grpcResp.Token)

	userId, err := parseUserId(grpcResp.UserId)
	if err != nil {
		h.log.Error("failed to parse user id", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Authentication failed")
		return
	}

	responseData := gen.SignInResponse{
		Email:  &grpcResp.Email,
		Exp:    &exp,
		Name:   &grpcResp.Name,
		Token:  &grpcResp.Token,
		UserId: &userId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}
