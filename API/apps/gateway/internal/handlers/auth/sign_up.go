package authhandler

import (
	"encoding/json"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/models"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req Credentials
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body")
		response.ErrorResponseJSON(w,r,http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validateCredentials(req); err != nil {
		h.log.Error("invalid credentials", "error", err)
		response.ErrorResponseJSON(w,r,http.StatusBadRequest, "invalid credentials")
		return
	}

	exists, err := h.storage.UserExists(ctx, req.Username, req.Email)
	if err != nil {
		h.log.Error("database error checking user existence")
		response.ErrorResponseJSON(w,r,http.StatusInternalServerError, "internal server error")
		return
	}

	if exists {
		h.log.Error("user with this username or email already exists")
		response.ErrorResponseJSON(w,r,http.StatusConflict, "user with this username or email already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.log.Error("failed to hash password")
		response.ErrorResponseJSON(w,r,http.StatusInternalServerError, "failed to process password")
		return
	}

	user := models.User{
		Name:     req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := h.storage.CreateUser(ctx, &user); err != nil {
		h.log.Error("failed to create user")
		response.ErrorResponseJSON(w,r,http.StatusInternalServerError, "failed to create user")
		return
	}

	response.SuccessResponse(w,r,http.StatusCreated, struct {
		Message string `json:"message"`
		UserID  int64  `json:"userID"`
	}{
		Message: "user created successfully",
		UserID:  user.ID,
	})

	return
}
