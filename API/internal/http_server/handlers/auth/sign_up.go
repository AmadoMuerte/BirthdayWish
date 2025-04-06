package authhandler

import (
	"encoding/json"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/lib/response"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/models"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req Credentials
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request body"))
		return
	}

	if err := validateCredentials(req); err != nil {
		h.log.Error("invalid credentials", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	exists, err := h.storage.UserExists(ctx, req.Username, req.Email)
	if err != nil {
		h.log.Error("database error checking user existence", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("internal server error"))
		return
	}

	if exists {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, response.Error("user with this username or email already exists"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.log.Error("failed to hash password", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to process password"))
		return
	}

	user := models.User{
		Name:     req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := h.storage.CreateUser(ctx, &user); err != nil {
		h.log.Error("failed to create user", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to create user"))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, struct {
		Message string `json:"message"`
		UserID  int64  `json:"userID"`
	}{
		Message: "user created successfully",
		UserID:  user.ID,
	})
	return
}
