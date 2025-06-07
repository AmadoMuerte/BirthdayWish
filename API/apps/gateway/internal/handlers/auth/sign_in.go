package authhandler

import (
	"encoding/json"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
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
		response.ErrorResponseJSON(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.storage.GetUserByUsername(ctx, credentials.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
		h.log.Error("Error get user by username", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	claims := jwt.NewClaims(jwt.Claims{UserID: user.ID})
	_, tokenString, err := h.tokenAuth.Encode(claims)
	if err != nil {
		h.log.Error("signup: failed to create user", "error", err, "userID", user.ID)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
	
	response.SuccessResponse(w, r, http.StatusOK, map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"token": tokenString,
	})

	return
}
