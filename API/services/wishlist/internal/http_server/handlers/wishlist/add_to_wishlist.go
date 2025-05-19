package wishlist

import (
	"encoding/json"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/lib/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/lib/response"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/models"
	"github.com/go-chi/render"
)

func (h *WishlistHandler) AddToWishlist(w http.ResponseWriter, r *http.Request) {
	var req models.Wishlist
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request body"))
		return
	}

	claims, err := jwt.GetClaims(ctx)
	if err != nil {
		h.log.Error("failed to get claims from token", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, response.Error("invalid token"))
		return
	}

	exists, err := h.storage.UserExistsByID(ctx, claims.UserID)
	if err != nil || !exists {
		h.log.Error("user does not exist", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("user does not exist"))
		return
	}
	req.UserID = claims.UserID

	if req.Name == "" {
		h.log.Error("invalid name")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid name"))
		return
	}

	if req.Link == "" {
		h.log.Error("invalid link")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid link"))
		return
	}

	if err := h.storage.AddToWishlist(ctx, req); err != nil {
		h.log.Error("failed to add wishlist item", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to add wishlist item"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, response.Success("wishlist item added successfully"))
}
