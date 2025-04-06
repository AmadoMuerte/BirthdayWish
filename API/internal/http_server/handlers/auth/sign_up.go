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
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to decode data"))
		return
	}

	err = validateCredentials(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error("failed to validate user")
		render.JSON(w, r, response.Error("invalid username or password"))
		return
	}

	user := models.User{
		Name:     req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error("failed to hash password")
		render.JSON(w, r, response.Error("a user with the same name already exists"))
		return
	}
	user.Password = string(hashedPassword)

	err = h.storage.CreateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		h.log.Error("error", err)
		render.JSON(w, r, response.Error("a user with the same name already exists"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
