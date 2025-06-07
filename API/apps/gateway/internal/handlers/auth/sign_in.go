package authhandler

import (
	"encoding/json"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		h.log.Error("error decode credentials", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, response.Error("Unauthorized"))
		return
	}

	user, err := h.storage.GetUserByUsername(ctx, credentials.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
		h.log.Error("Error get user by username", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, response.Error("Unauthorized"))
		return
	}

	claims := jwt.NewClaims(jwt.Claims{UserID: user.ID})
	_, tokenString, err := h.tokenAuth.Encode(claims)
	if err != nil {
		h.log.Error("signup: failed to create user", "error", err, "userID", user.ID)
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("Internal server error"))
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusOK)

	render.JSON(w, r, map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"token": tokenString,
	})

	return
}
